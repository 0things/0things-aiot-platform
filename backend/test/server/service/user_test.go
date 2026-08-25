package service_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/service"
	"aiot-backend/pkg/config"
	"aiot-backend/pkg/jwt"
	"aiot-backend/pkg/log"
	"aiot-backend/pkg/sid"
	mock_repository "aiot-backend/test/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var (
	logger *log.Logger
	j      *jwt.JWT
	sf     *sid.Sid
)

func TestMain(m *testing.M) {
	fmt.Println("begin")

	err := os.Setenv("APP_CONF", "../../../config/local.yml")
	if err != nil {
		panic(err)
	}

	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger = log.NewLog(conf)
	j = jwt.NewJwt(conf)
	sf = sid.NewSid()

	code := m.Run()
	fmt.Println("test end")

	os.Exit(code)
}

func TestUserService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockOrgRepo := mock_repository.NewMockOrganizationRepository(ctrl)
	mockOrgUserRepo := mock_repository.NewMockOrganizationUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)

	userService := service.NewUserService(srv, mockUserRepo, mockOrgRepo, mockOrgUserRepo)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Password: "password",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, nil)
	mockTm.EXPECT().Transaction(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
		return fn(ctx)
	})
	mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	mockOrgRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	mockOrgUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	err := userService.Register(ctx, req)

	assert.NoError(t, err)
}

func TestUserService_Register_UsernameExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockOrgRepo := mock_repository.NewMockOrganizationRepository(ctrl)
	mockOrgUserRepo := mock_repository.NewMockOrganizationUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo, mockOrgRepo, mockOrgUserRepo)

	ctx := context.Background()
	req := &v1.RegisterRequest{
		Password: "password",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(&model.User{}, nil)

	err := userService.Register(ctx, req)

	assert.Error(t, err)
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockOrgRepo := mock_repository.NewMockOrganizationRepository(ctrl)
	mockOrgUserRepo := mock_repository.NewMockOrganizationUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo, mockOrgRepo, mockOrgUserRepo)

	ctx := context.Background()
	req := &v1.LoginRequest{
		Email:    "xxx@gmail.com",
		Password: "password",
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Error("failed to hash password")
	}

	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(&model.User{
		UserId:   "user_123",
		Password: string(hashedPassword),
	}, nil)

	now := time.Now()
	mockOrgUserRepo.EXPECT().ListByUser(ctx, "user_123").Return([]*model.OrganizationUser{
		{OrganizationID: 1, UserID: "user_123", LastLoginAt: &now},
	}, nil)
	mockOrgUserRepo.EXPECT().UpdateLastLogin(ctx, "user_123", int64(1), gomock.Any()).Return(nil)

	token, err := userService.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockOrgRepo := mock_repository.NewMockOrganizationRepository(ctrl)
	mockOrgUserRepo := mock_repository.NewMockOrganizationUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo, mockOrgRepo, mockOrgUserRepo)

	ctx := context.Background()
	req := &v1.LoginRequest{
		Email:    "xxx@gmail.com",
		Password: "password",
	}

	mockUserRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, errors.New("user not found"))

	_, err := userService.Login(ctx, req)

	assert.Error(t, err)
}

func TestUserService_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockOrgRepo := mock_repository.NewMockOrganizationRepository(ctrl)
	mockOrgUserRepo := mock_repository.NewMockOrganizationUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo, mockOrgRepo, mockOrgUserRepo)

	ctx := context.Background()
	userId := "123"

	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(&model.User{
		UserId: userId,
		Email:  "test@example.com",
	}, nil)

	user, err := userService.GetProfile(ctx, userId)

	assert.NoError(t, err)
	assert.Equal(t, userId, user.UserId)
}

func TestUserService_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockOrgRepo := mock_repository.NewMockOrganizationRepository(ctrl)
	mockOrgUserRepo := mock_repository.NewMockOrganizationUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo, mockOrgRepo, mockOrgUserRepo)

	ctx := context.Background()
	userId := "123"
	req := &v1.UpdateProfileRequest{
		Nickname: "testuser",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(&model.User{
		UserId: userId,
		Email:  "old@example.com",
	}, nil)
	mockUserRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil)

	err := userService.UpdateProfile(ctx, userId, req)

	assert.NoError(t, err)
}

func TestUserService_UpdateProfile_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repository.NewMockUserRepository(ctrl)
	mockOrgRepo := mock_repository.NewMockOrganizationRepository(ctrl)
	mockOrgUserRepo := mock_repository.NewMockOrganizationUserRepository(ctrl)
	mockTm := mock_repository.NewMockTransaction(ctrl)
	srv := service.NewService(mockTm, logger, sf, j)
	userService := service.NewUserService(srv, mockUserRepo, mockOrgRepo, mockOrgUserRepo)

	ctx := context.Background()
	userId := "123"
	req := &v1.UpdateProfileRequest{
		Nickname: "testuser",
		Email:    "test@example.com",
	}

	mockUserRepo.EXPECT().GetByID(ctx, userId).Return(nil, errors.New("user not found"))

	err := userService.UpdateProfile(ctx, userId, req)

	assert.Error(t, err)
}
