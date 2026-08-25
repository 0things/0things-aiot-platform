package handler

import (
	"context"
	"testing"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/service"
	"aiot-backend/pkg/sid"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserOrgMembership_EndToEnd(t *testing.T) {
	// 1. 初始化内存数据库与迁移
	userDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, userDB.AutoMigrate(&model.User{}, &model.Organization{}, &model.OrganizationUser{}))

	deviceDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, deviceDB.AutoMigrate(&model.Product{}, &model.Device{}, &model.DeviceState{}))

	baseRepo := repository.NewRepository(logger, userDB)
	tm := repository.NewTransaction(baseRepo)
	userRepo := repository.NewUserRepository(baseRepo)
	orgRepo := repository.NewOrganizationRepository(baseRepo)
	orgUserRepo := repository.NewOrganizationUserRepository(baseRepo)

	svc := service.NewService(tm, logger, sid.NewSid(), jwt)
	userSvc := service.NewUserService(svc, userRepo, orgRepo, orgUserRepo)

	ctx := context.Background()

	// 2. 注册新用户
	regReq := &v1.RegisterRequest{
		Email:    "alice@example.com",
		Password: "password123",
	}
	require.NoError(t, userSvc.Register(ctx, regReq))

	// 3. 登录获取 Token
	loginReq := &v1.LoginRequest{
		Email:    "alice@example.com",
		Password: "password123",
	}
	token1, err := userSvc.Login(ctx, loginReq)
	require.NoError(t, err)
	require.NotEmpty(t, token1)

	claims1, err := jwt.ParseToken(token1)
	require.NoError(t, err)
	require.NotEmpty(t, claims1.UserId)
	require.NotZero(t, claims1.OrganizationID)
	org1ID := claims1.OrganizationID

	// 4. 查询组织列表
	orgs, err := userSvc.ListMyOrganizations(ctx, claims1.UserId)
	require.NoError(t, err)
	require.Len(t, orgs, 1)
	require.Equal(t, org1ID, orgs[0].Id)

	// 5. 新建第二个组织并建立关联
	org2 := &model.Organization{Id: 2, Name: "研发团队", Slug: "rd"}
	require.NoError(t, orgRepo.Create(ctx, org2))
	require.NoError(t, orgUserRepo.Create(ctx, &model.OrganizationUser{
		OrganizationID: 2,
		UserID:         claims1.UserId,
	}))

	// 6. 切换组织到 org2
	token2, err := userSvc.SwitchOrganization(ctx, claims1.UserId, 2)
	require.NoError(t, err)
	claims2, err := jwt.ParseToken(token2)
	require.NoError(t, err)
	require.Equal(t, int64(2), claims2.OrganizationID)

	// 7. 非归属组织切换被拒绝
	_, err = userSvc.SwitchOrganization(ctx, claims1.UserId, 999)
	require.Error(t, err)
	require.Equal(t, v1.ErrForbidden, err)
}
