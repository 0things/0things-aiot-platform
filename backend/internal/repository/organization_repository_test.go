package repository

import (
	"context"
	"testing"
	"time"

	"aiot-backend/internal/model"

	"github.com/stretchr/testify/require"
)

func TestOrganizationRepository(t *testing.T) {
	store := newRepositoryTestDB(t, &model.Organization{})
	baseRepo := &Repository{db: store}
	orgRepo := NewOrganizationRepository(baseRepo)
	ctx := context.Background()

	org1 := &model.Organization{Name: "Org 1", Slug: "org-1"}
	require.NoError(t, orgRepo.Create(ctx, org1))
	require.NotZero(t, org1.Id)

	org2 := &model.Organization{Name: "Org 2", Slug: "org-2"}
	require.NoError(t, orgRepo.Create(ctx, org2))

	fetched, err := orgRepo.GetByID(ctx, org1.Id)
	require.NoError(t, err)
	require.Equal(t, "Org 1", fetched.Name)

	list, err := orgRepo.ListByIDs(ctx, []int64{org1.Id, org2.Id})
	require.NoError(t, err)
	require.Len(t, list, 2)

	emptyList, err := orgRepo.ListByIDs(ctx, []int64{})
	require.NoError(t, err)
	require.Empty(t, emptyList)
}

func TestOrganizationUserRepository(t *testing.T) {
	store := newRepositoryTestDB(t, &model.OrganizationUser{})
	baseRepo := &Repository{db: store}
	orgUserRepo := NewOrganizationUserRepository(baseRepo)
	ctx := context.Background()

	// 1. Create and Member check
	ou1 := &model.OrganizationUser{
		OrganizationID: 1,
		UserID:         "user_001",
	}
	require.NoError(t, orgUserRepo.Create(ctx, ou1))

	isMember, err := orgUserRepo.IsMember(ctx, "user_001", 1)
	require.NoError(t, err)
	require.True(t, isMember)

	isMember, err = orgUserRepo.IsMember(ctx, "user_001", 2)
	require.NoError(t, err)
	require.False(t, isMember)

	// 2. Unique constraint check
	ouDuplicate := &model.OrganizationUser{
		OrganizationID: 1,
		UserID:         "user_001",
	}
	require.Error(t, orgUserRepo.Create(ctx, ouDuplicate))

	// 3. Multiple orgs
	ou2 := &model.OrganizationUser{
		OrganizationID: 2,
		UserID:         "user_001",
	}
	require.NoError(t, orgUserRepo.Create(ctx, ou2))

	orgIDs, err := orgUserRepo.ListOrgIDsByUser(ctx, "user_001")
	require.NoError(t, err)
	require.ElementsMatch(t, []int64{1, 2}, orgIDs)

	// 4. Update last_login_at
	now := time.Now().Truncate(time.Second)
	require.NoError(t, orgUserRepo.UpdateLastLogin(ctx, "user_001", 2, now))

	fetched, err := orgUserRepo.GetByUserAndOrg(ctx, "user_001", 2)
	require.NoError(t, err)
	require.NotNil(t, fetched.LastLoginAt)
	require.Equal(t, now.Unix(), fetched.LastLoginAt.Unix())
}
