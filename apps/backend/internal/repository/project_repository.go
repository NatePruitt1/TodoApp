package repository

import (
	"backend/internal/models"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository interface {
	Create(project models.Project) error
	Delete(projectId uuid.UUID) error
	GetByID(projectId uuid.UUID) (models.Project, error)
	GetAllUserProjects(userId uuid.UUID) ([]models.Project, error)
}

type ProjectRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepositoryImpl {
	return &ProjectRepositoryImpl{
		db: db,
	}
}

func (ur *ProjectRepositoryImpl) Create(project models.Project) error {
	const q = `
		INSERT INTO projects (id, project_name, description, owner_id)
		VALUES ($1, $2, $3, $4);
	`

	_, err := ur.db.Exec(context.Background(), q,
		project.ID,
		project.Name,
		project.Description,
		project.OwnerID,
	)

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

func (ur *ProjectRepositoryImpl) GetByID(id uuid.UUID) (models.Project, error) {
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

	return project, err
}

func (ur *ProjectRepositoryImpl) GetAllUserProjects(owner_id uuid.UUID) ([]models.Project, error) {
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

	projects := make([]models.Project, 0)
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

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
