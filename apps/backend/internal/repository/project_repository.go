package repository

import (
	"backend/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrProjectNotFound     = errors.New("project not found.")
	ErrCategoryNotFound    = errors.New("category not found.")
	ErrCardNotFound        = errors.New("card not found.")
	ErrNotResourceOwner    = errors.New("user does not own resource.")
	ErrProjectNameConflict = errors.New("project name already taken for user.")
	ErrInvalidOwner        = errors.New("owner does not exist.")
)

type ProjectRepository interface {
	Delete(projectId uuid.UUID) error
	GetByID(projectId uuid.UUID) (*models.Project, error)
	GetAllUserProjects(userId uuid.UUID) ([]*models.Project, error)
	GetAggregate(project uuid.UUID) (*models.Project, error)
	Save(project *models.Project) error
	CheckProjectOwner(projectId, userId uuid.UUID) error

	// Category CRUD
	DeleteCategory(catId uuid.UUID) error
	GetCategoryByID(catId uuid.UUID) (*models.Category, error)
	GetCategoriesForProject(projectId uuid.UUID) ([]*models.Category, error)
	SaveCategory(category *models.Category) error
	CheckCategoryOwner(catId, userId uuid.UUID) error

	// Card CRUD
	DeleteCard(cardId uuid.UUID) error
	GetCardByID(cardId uuid.UUID) (*models.Card, error)
	GetCardsForCategory(categoryId uuid.UUID) ([]*models.Card, error)
	SaveCard(card *models.Card) error
	CheckCardOwner(cardId, userId uuid.UUID) error
}

type ProjectRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) ProjectRepository {
	return &ProjectRepositoryImpl{
		db: db,
	}
}

//--------------------- Helpers ---------------------

func (ur *ProjectRepositoryImpl) queryOne(ctx context.Context, notFoundError error, query string, dest []interface{}, args ...interface{}) error {
	err := ur.db.QueryRow(ctx, query, args...).Scan(dest...)
	if err == pgx.ErrNoRows {
		return notFoundError
	}
	return err
}

// Pass in pool instead of reciever function due to golang type limitations.
func queryMany[T any](ctx context.Context, pool *pgxpool.Pool, query string, notFoundError error, args []interface{}, scan func(pgx.Rows) (*T, error), funcName string) ([]*T, error) {
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, notFoundError
		} else {
			return nil, fmt.Errorf("%s :%w", funcName, err)
		}
	}

	defer rows.Close()

	results := make([]*T, 0)
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", funcName, err)
		}
		results = append(results, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return results, nil
}

func (ur *ProjectRepositoryImpl) upsert(ctx context.Context, insertQ, updateQ, funcName string, notFoundError error, args ...interface{}) error {
	tag, err := ur.db.Exec(ctx, insertQ, args...)
	if pgErr := (*pgconn.PgError)(nil); errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		utag, err := ur.db.Exec(ctx, updateQ, args...)
		if err != nil {
			return fmt.Errorf("%s: update after conflict: %w", funcName, err)
		}
		if utag.RowsAffected() == 0 {
			return notFoundError
		}
	} else if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	if tag.RowsAffected() == 0 {
		return notFoundError
	}

	return err
}

func (ur *ProjectRepositoryImpl) checkOwner(ctx context.Context, query string, resourceID, userID uuid.UUID) error {
	var ownerId uuid.UUID
	err := ur.db.QueryRow(ctx, query, resourceID).Scan(&ownerId)
	if err != nil {
		return err
	}
	if ownerId != userID {
		return ErrNotResourceOwner
	}
	return nil
}

func (ur *ProjectRepositoryImpl) deleteByID(ctx context.Context, query string, id uuid.UUID, notFoundError error) error {
	tag, err := ur.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return notFoundError
	}
	return nil
}

//------------ Interface IMPL --------------------

func (ur *ProjectRepositoryImpl) Save(project *models.Project) error {
	const insertQ = `
		INSERT INTO projects (id, project_name, description, owner_id)
		VALUES ($1, $2, $3, $4);
	`

	const updateQ = `
		UPDATE projects
		SET project_name = $2, description = $3, owner_id = $4
		WHERE id = $1
	`

	return ur.upsert(context.Background(), insertQ, updateQ, "SaveProject", ErrProjectNotFound, project.ID, project.Name, project.Description, project.OwnerID)
}

func (ur *ProjectRepositoryImpl) Delete(id uuid.UUID) error {
	const q = `
		DELETE FROM projects
		WHERE id = $1
		
	`

	return ur.deleteByID(context.Background(), q, id, ErrProjectNotFound)
}

func (ur *ProjectRepositoryImpl) GetByID(id uuid.UUID) (*models.Project, error) {
	const q = `
		SELECT id, project_name, description, owner_id 
		FROM projects
		WHERE id = $1
	`
	var project models.Project
	err := ur.queryOne(context.Background(), ErrProjectNotFound, q,
		[]interface{}{&project.ID, &project.Name, &project.Description, &project.OwnerID},
		id,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (ur *ProjectRepositoryImpl) GetAllUserProjects(owner_id uuid.UUID) ([]*models.Project, error) {
	const q = `
		SELECT id, project_name, description, owner_id 
		FROM projects
		WHERE owner_id = $1
	`

	return queryMany(
		context.Background(),
		ur.db,
		q,
		ErrProjectNotFound,
		[]interface{}{owner_id},
		func(r pgx.Rows) (*models.Project, error) {
			var p models.Project
			err := r.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID)
			return &p, err
		},
		"GetAllUserProjects",
	)
}

func (ur *ProjectRepositoryImpl) GetAggregate(projectId uuid.UUID) (*models.Project, error) {
	//get project
	project, err := ur.GetByID(projectId)
	if err != nil {
		return nil, err
	}

	categories, err := ur.GetCategoriesForProject(project.ID)
	if err != nil {
		return nil, err
	}

	for c := range categories {
		category := categories[c]
		cards, err := ur.GetCardsForCategory(category.ID)
		if err != nil {
			return nil, err
		}

		category.Cards = cards
	}

	project.Categories = categories
	return project, nil
}

func (ur *ProjectRepositoryImpl) DeleteCategory(catId uuid.UUID) error {
	const q = `DELETE FROM categories WHERE id = $1`
	return ur.deleteByID(context.Background(), q, catId, ErrCategoryNotFound)
}

func (ur *ProjectRepositoryImpl) GetCategoryByID(catId uuid.UUID) (*models.Category, error) {
	const q = `SELECT id, project_id, category_name, index FROM categories WHERE id = $1`

	var category models.Category
	err := ur.queryOne(
		context.Background(),
		ErrCategoryNotFound,
		q,
		[]interface{}{&category.ID, &category.ProjectID, &category.Name, &category.Index},
		catId,
	)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (ur *ProjectRepositoryImpl) GetCategoriesForProject(projectId uuid.UUID) ([]*models.Category, error) {
	const q = `SELECT id, project_id, category_name, index FROM categories WHERE project_id = $1 ORDER BY index ASC`

	return queryMany(
		context.Background(),
		ur.db,
		q,
		ErrCategoryNotFound,
		[]interface{}{projectId},
		func(r pgx.Rows) (*models.Category, error) {
			var c models.Category
			err := r.Scan(&c.ID, &c.ProjectID, &c.Name, &c.Index)
			return &c, err
		},
		"GetCategoriesForProject",
	)
}

func (ur *ProjectRepositoryImpl) SaveCategory(category *models.Category) error {
	const insertQ = `INSERT INTO categories (id, project_id, category_name, index) VALUES ($1, $2, $3, $4)`
	const updateQ = `UPDATE categories SET project_id = $2, category_name = $3, index = $4 WHERE id = $1`
	return ur.upsert(context.Background(), insertQ, updateQ, "SaveCategory", ErrCategoryNotFound, category.ID, category.ProjectID, category.Name, category.Index)
}

func (ur *ProjectRepositoryImpl) CheckProjectOwner(projectId, userId uuid.UUID) error {
	const q = `SELECT owner_id FROM projects WHERE id = $1`
	return ur.checkOwner(context.Background(), q, projectId, userId)
}

func (ur *ProjectRepositoryImpl) CheckCategoryOwner(catId, userId uuid.UUID) error {
	const q = `SELECT p.owner_id FROM categories c JOIN projects p ON p.id = c.project_id WHERE c.id = $1`
	return ur.checkOwner(context.Background(), q, catId, userId)
}

func (ur *ProjectRepositoryImpl) DeleteCard(cardId uuid.UUID) error {
	const q = `DELETE FROM cards WHERE id = $1`
	return ur.deleteByID(context.Background(), q, cardId, ErrCardNotFound)
}

func (ur *ProjectRepositoryImpl) GetCardByID(cardId uuid.UUID) (*models.Card, error) {
	const q = `SELECT id, category_id, title, content, finished FROM cards WHERE id = $1`

	var card models.Card
	err := ur.queryOne(
		context.Background(),
		ErrCardNotFound,
		q,
		[]interface{}{&card.ID, &card.CategoryID, &card.Title, &card.Content, &card.Finished},
		cardId,
	)

	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (ur *ProjectRepositoryImpl) GetCardsForCategory(categoryId uuid.UUID) ([]*models.Card, error) {
	const q = `SELECT id, title, content, category_id, finished FROM cards WHERE category_id = $1`

	return queryMany(
		context.Background(),
		ur.db,
		q,
		ErrCategoryNotFound,
		[]interface{}{categoryId},
		func(r pgx.Rows) (*models.Card, error) {
			var c models.Card
			err := r.Scan(&c.ID, &c.Title, &c.Content, &c.CategoryID, &c.Finished)
			return &c, err
		},
		"GetCardsForCategory",
	)
}

func (ur *ProjectRepositoryImpl) SaveCard(card *models.Card) error {
	const insertQ = `INSERT INTO cards (id, category_id, title, content, finished) VALUES ($1, $2, $3, $4, $5)`
	const updateQ = `UPDATE cards SET category_id = $2, title = $3, content = $4, finished = $5 WHERE id = $1`
	return ur.upsert(context.Background(), insertQ, updateQ, "SaveCard", ErrCardNotFound, card.ID, card.CategoryID, card.Title, card.Content, card.Finished)
}

func (ur *ProjectRepositoryImpl) CheckCardOwner(cardId, userId uuid.UUID) error {
	const q = `SELECT p.owner_id FROM cards ca JOIN categories c ON c.id = ca.category_id JOIN projects p ON p.id = c.project_id WHERE ca.id = $1`
	return ur.checkOwner(context.Background(), q, cardId, userId)
}
