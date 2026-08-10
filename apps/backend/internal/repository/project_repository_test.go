package repository

import (
	"backend/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectRepo_SaveNoError(t *testing.T) {
	pool, m := setTestDB(t)
	defer func() {
		require.NoError(t, m.Down())
	}()
	repo := NewProjectRepository(pool)

	user := makeTestUser(t, pool, "username", "abcd")

	project := models.Project{
		ID:          uuid.New(),
		OwnerID:     user.ID,
		Name:        "new project",
		Description: "project description",
		Categories:  nil,
	}

	require.NoError(t, repo.Save(&project))
}

func TestProjectRepo_GetByID(t *testing.T) {
	pool, m := setTestDB(t)
	defer func() {
		require.NoError(t, m.Down())
	}()
	repo := NewProjectRepository(pool)

	user := makeTestUser(t, pool, "username", "abcd")

	project := models.Project{
		ID:          uuid.New(),
		OwnerID:     user.ID,
		Name:        "new project",
		Description: "project description",
		Categories:  nil,
	}

	require.NoError(t, repo.Save(&project))

	got, err := repo.GetByID(project.ID)
	require.NoError(t, err)

	assert.Equal(t, project, *got)
}

func TestProjectRepo_GetByIDNotExist(t *testing.T) {
	pool, m := setTestDB(t)
	defer func() {
		require.NoError(t, m.Down())
	}()
	repo := NewProjectRepository(pool)
	proj, err := repo.GetByID(uuid.New())
	assert.Nil(t, proj)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrProjectNotFound)
}

func TestProjectRepo_GetAllUserProjects(t *testing.T) {
	pool, m := setTestDB(t)
	defer func() {
		require.NoError(t, m.Down())
	}()
	repo := NewProjectRepository(pool)

	user := makeTestUser(t, pool, "username", "password")
	user2 := makeTestUser(t, pool, "user2", "abcd")
	proj1 := makeTestProject(t, pool, "p1", user, "desc")
	proj2 := makeTestProject(t, pool, "p2", user, "desc")
	proj3 := makeTestProject(t, pool, "p3", user, "desc")

	proj4 := makeTestProject(t, pool, "p1", user2, "desc")

	projects, err := repo.GetAllUserProjects(user.ID)
	require.NoError(t, err)

	projectsUnwraped := make([]models.Project, len(projects))
	for p := range projects {
		projectsUnwraped[p] = *projects[p]
	}

	assert.Contains(t, projectsUnwraped, proj1)
	assert.Contains(t, projectsUnwraped, proj2)
	assert.Contains(t, projectsUnwraped, proj3)
	assert.NotContains(t, projectsUnwraped, proj4)
}

func TestProjectRepo_CheckProjectOwner(t *testing.T) {
	pool, m := setTestDB(t)
	defer func() {
		require.NoError(t, m.Down())
	}()
	repo := NewProjectRepository(pool)

	user := makeTestUser(t, pool, "user", "abcd")
	proj1 := makeTestProject(t, pool, "p1", user, "desc")

	err := repo.CheckProjectOwner(proj1.ID, user.ID)
	require.NoError(t, err)
}
