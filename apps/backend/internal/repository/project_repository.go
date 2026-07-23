package repository

import (
	"backend/internal/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
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

	_, err := ur.db.Exec(context.Background(), insertQ,
		project.ID,
		project.Name,
		project.Description,
		project.OwnerID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				_, err = ur.db.Exec(context.Background(), updateQ,
					project.ID,
					project.Name,
					project.Description,
					project.OwnerID,
				)
			}
		}
	}

	return err
}

func (ur *ProjectRepositoryImpl) Delete(id uuid.UUID) error {
	const q = `
		DELETE FROM projects
		WHERE id = $1
		
	`

	_, err := ur.db.Exec(context.Background(), q, id)
	return err
}

func (ur *ProjectRepositoryImpl) GetByID(id uuid.UUID) (*models.Project, error) {
	const q = `
		SELECT id, project_name, description, owner_id 
		FROM projects
		WHERE id = $1
	`

	var project models.Project
	err := ur.db.QueryRow(context.Background(), q, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.OwnerID,
	)

	return &project, err
}

func (ur *ProjectRepositoryImpl) GetAllUserProjects(owner_id uuid.UUID) ([]*models.Project, error) {
	const q = `
		SELECT id, project_name, description, owner_id 
		FROM projects
		WHERE owner_id = $1
	`

	rows, err := ur.db.Query(context.Background(), q, owner_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]*models.Project, 0)
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.OwnerID,
		); err != nil {
			return nil, err
		}

		projects = append(projects, &project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
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
	const q = `
		DELETE FROM categories
		WHERE id = $1
	`

	_, err := ur.db.Exec(context.Background(), q, catId)
	return err
}

func (ur *ProjectRepositoryImpl) GetCategoryByID(catId uuid.UUID) (*models.Category, error) {
	const q = `
		SELECT id, project_id, category_name, index
		FROM categories
		WHERE id = $1
	`

	var category models.Category
	err := ur.db.QueryRow(context.Background(), q, catId).Scan(
		&category.ID,
		&category.ProjectID,
		&category.Name,
		&category.Index,
	)

	return &category, err
}

func (ur *ProjectRepositoryImpl) GetCategoriesForProject(projectId uuid.UUID) ([]*models.Category, error) {
	const q = `
		SELECT id, project_id, category_name, index
		FROM categories
		WHERE project_id = $1
		ORDER BY index ASC
	`

	rows, err := ur.db.Query(context.Background(), q, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*models.Category, 0)
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(
			&category.ID,
			&category.ProjectID,
			&category.Name,
			&category.Index,
		); err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (ur *ProjectRepositoryImpl) SaveCategory(category *models.Category) error {
	const insertQ = `
		INSERT INTO categories (id, project_id, category_name, index)
		VALUES ($1, $2, $3, $4)
	`

	const updateQ = `
		UPDATE categories
		SET project_id = $2, category_name = $3, index = $4
		WHERE id = $1
	`

	_, err := ur.db.Exec(context.Background(), insertQ,
		category.ID,
		category.ProjectID,
		category.Name,
		category.Index,
	)

	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			_, err = ur.db.Exec(context.Background(), updateQ,
				category.ID,
				category.ProjectID,
				category.Name,
				category.Index,
			)
		}
	}

	return err
}

func (ur *ProjectRepositoryImpl) CheckProjectOwner(projectId, userId uuid.UUID) error {
	const q = `
		SELECT owner_id FROM projects WHERE id = $1
	`

	var ownerId uuid.UUID
	err := ur.db.QueryRow(context.Background(), q, projectId).Scan(&ownerId)

	if err != nil {
		return err
	}

	if ownerId != userId {
		return errors.New("User does not own project.")
	}

	return nil
}

func (ur *ProjectRepositoryImpl) CheckCategoryOwner(catId, userId uuid.UUID) error {
	const q = `
		SELECT p.owner_id FROM categories c JOIN projects p ON p.id = c.project_id where c.id = $1
	`

	var ownerId uuid.UUID
	err := ur.db.QueryRow(context.Background(), q, catId).Scan(&ownerId)

	if err != nil {
		return err
	}

	if ownerId != userId {
		return errors.New("User does not own category.")
	}

	return nil
}

func (ur *ProjectRepositoryImpl) DeleteCard(cardId uuid.UUID) error {
	const q = `
		DELETE FROM cards
		WHERE id = $1
	`

	_, err := ur.db.Exec(context.Background(), q, cardId)
	return err
}

func (ur *ProjectRepositoryImpl) GetCardByID(cardId uuid.UUID) (*models.Card, error) {
	const q = `
		SELECT id, category_id, title, content, finished
		FROM cards
		WHERE id = $1
	`

	var card models.Card
	err := ur.db.QueryRow(context.Background(), q, cardId).Scan(
		&card.ID,
		&card.CategoryID,
		&card.Title,
		&card.Content,
		&card.Finished,
	)

	return &card, err
}

func (ur *ProjectRepositoryImpl) GetCardsForCategory(categoryId uuid.UUID) ([]*models.Card, error) {
	const q = `
		SELECT id, title, content, category_id, finished
		FROM cards
		WHERE category_id = $1
	`

	rows, err := ur.db.Query(context.Background(), q, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make([]*models.Card, 0)
	for rows.Next() {
		var card models.Card
		if err := rows.Scan(
			&card.ID,
			&card.Title,
			&card.Content,
			&card.CategoryID,
			&card.Finished,
		); err != nil {
			return nil, err
		}

		cards = append(cards, &card)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func (ur *ProjectRepositoryImpl) SaveCard(card *models.Card) error {
	const insertQ = `
		INSERT INTO cards (id, category_id, title, content, finished)
		VALUES ($1, $2, $3, $4, $5)
	`

	const updateQ = `
		UPDATE cards
		SET category_id = $2, title = $3, content = $4, finished = $5
		WHERE id = $1
	`

	_, err := ur.db.Exec(context.Background(), insertQ,
		card.ID,
		card.CategoryID,
		card.Title,
		card.Content,
		card.Finished,
	)

	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			_, err = ur.db.Exec(context.Background(), updateQ,
				card.ID,
				card.CategoryID,
				card.Title,
				card.Content,
				card.Finished,
			)
		}
	}

	return err
}

func (ur *ProjectRepositoryImpl) CheckCardOwner(cardId, userId uuid.UUID) error {
	const q = `
		SELECT p.owner_id FROM cards ca
		JOIN categories c ON c.id = ca.category_id
		JOIN projects p ON p.id = c.project_id
		WHERE ca.id = $1
	`

	var ownerId uuid.UUID
	err := ur.db.QueryRow(context.Background(), q, cardId).Scan(&ownerId)

	if err != nil {
		return err
	}

	if ownerId != userId {
		return errors.New("User does not own card.")
	}

	return nil
}
