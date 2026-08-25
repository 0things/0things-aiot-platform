package repository_test

import (
	"context"
	"testing"
	"time"

	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupSQLiteOrgRepos(t *testing.T) (repository.OrganizationRepository, repository.OrganizationUserRepository) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Organization{}, &model.OrganizationUser{}))
	repo := repository.NewRepository(nil, db)
	return repository.NewOrganizationRepository(repo), repository.NewOrganizationUserRepository(repo)
}

func TestOrganizationRepository_CRUD(t *testing.T) {
	orgRepo, _ := setupSQLiteOrgRepos(t)
	ctx := context.Background()

	// 1. Create
	org1 := &model.Organization{Name: "Org 1", Slug: "org-1"}
	err := orgRepo.Create(ctx, org1)
	require.NoError(t, err)
	assert.NotZero(t, org1.Id)

	org2 := &model.Organization{Name: "Org 2", Slug: "org-2"}
	err = orgRepo.Create(ctx, org2)
	require.NoError(t, err)

	// 2. GetByID
	found, err := orgRepo.GetByID(ctx, org1.Id)
	require.NoError(t, err)
	assert.Equal(t, "Org 1", found.Name)

	// 3. GetByID NotFound
	_, err = orgRepo.GetByID(ctx, 9999)
	assert.Error(t, err)

	// 4. ListByIDs
	list, err := orgRepo.ListByIDs(ctx, []int64{org1.Id, org2.Id})
	require.NoError(t, err)
	assert.Len(t, list, 2)

	// 5. ListByIDs Empty
	emptyList, err := orgRepo.ListByIDs(ctx, []int64{})
	require.NoError(t, err)
	assert.Empty(t, emptyList)
}

func TestOrganizationUserRepository_CRUD(t *testing.T) {
	_, orgUserRepo := setupSQLiteOrgRepos(t)
	ctx := context.Background()

	// 1. Create
	now := time.Now()
	ou1 := &model.OrganizationUser{
		UserID:         "u1",
		OrganizationID: 1,
		LastLoginAt:    &now,
	}
	err := orgUserRepo.Create(ctx, ou1)
	require.NoError(t, err)

	ou2 := &model.OrganizationUser{
		UserID:         "u1",
		OrganizationID: 2,
	}
	require.NoError(t, orgUserRepo.Create(ctx, ou2))

	// 2. ListByUser
	items, err := orgUserRepo.ListByUser(ctx, "u1")
	require.NoError(t, err)
	assert.Len(t, items, 2)

	// 3. ListOrgIDsByUser
	orgIDs, err := orgUserRepo.ListOrgIDsByUser(ctx, "u1")
	require.NoError(t, err)
	assert.ElementsMatch(t, []int64{1, 2}, orgIDs)

	// 4. IsMember
	isMember, err := orgUserRepo.IsMember(ctx, "u1", 1)
	require.NoError(t, err)
	assert.True(t, isMember)

	isNotMember, err := orgUserRepo.IsMember(ctx, "u1", 999)
	require.NoError(t, err)
	assert.False(t, isNotMember)

	// 5. GetByUserAndOrg
	item, err := orgUserRepo.GetByUserAndOrg(ctx, "u1", 1)
	require.NoError(t, err)
	require.NotNil(t, item)
	assert.Equal(t, int64(1), item.OrganizationID)

	itemNotFound, err := orgUserRepo.GetByUserAndOrg(ctx, "u1", 999)
	require.NoError(t, err)
	assert.Nil(t, itemNotFound)

	// 6. UpdateLastLogin
	newLogin := time.Now().Add(time.Hour)
	err = orgUserRepo.UpdateLastLogin(ctx, "u1", 1, newLogin)
	require.NoError(t, err)

	itemUpdated, err := orgUserRepo.GetByUserAndOrg(ctx, "u1", 1)
	require.NoError(t, err)
	require.NotNil(t, itemUpdated)
	assert.NotNil(t, itemUpdated.LastLoginAt)
}
