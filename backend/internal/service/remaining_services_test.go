package service

import (
	"context"
	"testing"
	"time"

	messageParserV1 "aiot-backend/api/message_parser/v1"
	v1 "aiot-backend/api/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"
	"aiot-backend/internal/tenant"

	"aiot-backend/pkg/jwt"
	"aiot-backend/pkg/log"
	"aiot-backend/pkg/sid"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func iotDB(db *gorm.DB) *repository.IoTDB { return &repository.IoTDB{DB: db} }

func TestDeviceEventService_RecordAndList(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Product{}, &model.Device{}, &model.DeviceState{}, &model.DeviceEvent{}))
	require.NoError(t, db.Create(&model.Product{ID: 1, ProductKey: "P001", Name: "p", OrganizationID: 1}).Error)
	require.NoError(t, db.Create(&model.Device{ID: 1, DeviceKey: "D001", Name: "d", ProductID: 1, OrganizationID: 1}).Error)
	require.NoError(t, db.Create(&model.DeviceState{ID: 1, DeviceKey: "D001", State: "online"}).Error)

	svc := NewDeviceEventService(
		repository.NewDeviceEventRepository(iotDB(db)),
		repository.NewDeviceRepository(iotDB(db), &repository.IoTRedis{}),
	)
	ctx := context.Background()

	require.NoError(t, svc.Record(ctx, "P001", "D001", "online", 0, map[string]any{"a": 1}))

	// mismatched product key
	require.Error(t, svc.Record(ctx, "P999", "D001", "online", 0, nil))
	// missing fields
	require.Error(t, svc.Record(ctx, "", "D001", "online", 0, nil))

	list, n, err := svc.List(ctx, 1, 10, "", "D001", "", nil, nil)
	require.NoError(t, err)
	require.Equal(t, int64(1), n)
	require.NotEmpty(t, list)
}

func TestProductTSLService_CRUD(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Product{}, &model.ProductTSL{}))
	require.NoError(t, db.Create(&model.Product{ID: 1, ProductKey: "P001", Name: "p", OrganizationID: 1}).Error)

	svc := NewProductTSLService(
		repository.NewProductRepository(iotDB(db)),
		repository.NewProductTSLRepository(iotDB(db)),
	)
	ctx := context.Background()

	_, err = svc.Get(ctx, "P001")
	require.Error(t, err)

	require.NoError(t, svc.Upsert(ctx, "P001", `{"tsl":1}`))
	got, err := svc.Get(ctx, "P001")
	require.NoError(t, err)
	require.Equal(t, `{"tsl":1}`, got.TSL)

	require.NoError(t, svc.Upsert(ctx, "P001", `{"tsl":2}`))
	got, _ = svc.Get(ctx, "P001")
	require.Equal(t, `{"tsl":2}`, got.TSL)

	require.NoError(t, svc.Delete(ctx, "P001"))
}

func TestProductMessageParserService(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Product{}, &model.ProductMessageParser{}))
	require.NoError(t, db.Create(&model.Product{ID: 1, ProductKey: "P001", Name: "p", OrganizationID: 1}).Error)

	svc := NewProductMessageParserService(
		repository.NewProductRepository(iotDB(db)),
		repository.NewProductMessageParserRepository(iotDB(db)),
	)
	ctx := context.Background()

	parser, isDefault, err := svc.Get(ctx, "P001")
	require.NoError(t, err)
	require.True(t, isDefault)
	require.NotEmpty(t, parser.Script)

	saved, err := svc.Save(ctx, "P001", messageParserV1.LanguageJavaScriptES5, defaultProductMessageParserScript)
	require.NoError(t, err)
	require.Equal(t, messageParserV1.LanguageJavaScriptES5, saved.Language)

	// unsupported language
	_, err = svc.Save(ctx, "P001", "python", defaultProductMessageParserScript)
	require.Error(t, err)

	// Execute should run the default script without panicking.
	_, _ = svc.Execute(ctx, "P001", messageParserV1.ExecuteProductMessageParserRequest{
		Mode:    "custom",
		Topic:   "t",
		RawData: "00",
	})
}

func TestUserService_CRUD(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Organization{}, &model.OrganizationUser{}))

	logger := log.NewLog(viper.New())
	repo := repository.NewRepository(logger, db)
	tm := repository.NewTransaction(repo)
	jwtInstance := jwt.NewJwt(viper.New())
	svc := NewService(tm, logger, sid.NewSid(), jwtInstance)
	orgRepo := repository.NewOrganizationRepository(repo)
	orgUserRepo := repository.NewOrganizationUserRepository(repo)
	userSvc := NewUserService(svc, repository.NewUserRepository(repo), orgRepo, orgUserRepo)
	ctx := context.Background()

	// 1. 注册成功
	require.NoError(t, userSvc.Register(ctx, &v1.RegisterRequest{Email: "u@e.com", Password: "secret"}))
	// 重复邮箱注册失败
	require.Error(t, userSvc.Register(ctx, &v1.RegisterRequest{Email: "u@e.com", Password: "secret"}))

	var u model.User
	require.NoError(t, db.Where("email = ?", "u@e.com").First(&u).Error)
	require.NotEmpty(t, u.UserId)

	// 检查自动创建的 organization 和 organization_users
	var orgUsers []model.OrganizationUser
	require.NoError(t, db.Where("user_id = ?", u.UserId).Find(&orgUsers).Error)
	require.Len(t, orgUsers, 1)
	initialOrgID := orgUsers[0].OrganizationID

	// 2. 登录 - 首次登录（last_login_at 为 NULL，回退最小 org_id）
	token, err := userSvc.Login(ctx, &v1.LoginRequest{Email: "u@e.com", Password: "secret"})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := jwtInstance.ParseToken(token)
	require.NoError(t, err)
	require.Equal(t, u.UserId, claims.UserId)
	require.Equal(t, initialOrgID, claims.OrganizationID)

	// 检查 last_login_at 是否更新
	var ouUpdated model.OrganizationUser
	require.NoError(t, db.Where("user_id = ? AND organization_id = ?", u.UserId, initialOrgID).First(&ouUpdated).Error)
	require.NotNil(t, ouUpdated.LastLoginAt)

	// 密码错误
	_, err = userSvc.Login(ctx, &v1.LoginRequest{Email: "u@e.com", Password: "wrong"})
	require.Error(t, err)

	// 3. 多组织场景：再手动加入两个组织（ID: 10, 20）
	org10 := &model.Organization{Id: 10, Name: "Org 10", Slug: "org-10"}
	org20 := &model.Organization{Id: 20, Name: "Org 20", Slug: "org-20"}
	require.NoError(t, orgRepo.Create(ctx, org10))
	require.NoError(t, orgRepo.Create(ctx, org20))

	require.NoError(t, orgUserRepo.Create(ctx, &model.OrganizationUser{OrganizationID: 10, UserID: u.UserId}))
	require.NoError(t, orgUserRepo.Create(ctx, &model.OrganizationUser{OrganizationID: 20, UserID: u.UserId}))

	// 给 org 20 设置最新的 last_login_at
	latestTime := time.Now().Add(time.Hour)
	require.NoError(t, orgUserRepo.UpdateLastLogin(ctx, u.UserId, 20, latestTime))

	// 再次登录，应当选中 org 20
	token2, err := userSvc.Login(ctx, &v1.LoginRequest{Email: "u@e.com", Password: "secret"})
	require.NoError(t, err)
	claims2, err := jwtInstance.ParseToken(token2)
	require.NoError(t, err)
	require.Equal(t, int64(20), claims2.OrganizationID)

	// 4. ListMyOrganizations (上下文当前组织为 20)
	ctxWithTenant := tenant.WithTenant(ctx, 20)
	orgList, err := userSvc.ListMyOrganizations(ctxWithTenant, u.UserId)
	require.NoError(t, err)
	require.Len(t, orgList, 3)

	var currentCount int
	for _, o := range orgList {
		if o.IsCurrent {
			currentCount++
			require.Equal(t, int64(20), o.Id)
		}
	}
	require.Equal(t, 1, currentCount)

	// 5. SwitchOrganization
	// 切换到已属组织 10
	switchToken, err := userSvc.SwitchOrganization(ctx, u.UserId, 10)
	require.NoError(t, err)
	switchClaims, err := jwtInstance.ParseToken(switchToken)
	require.NoError(t, err)
	require.Equal(t, int64(10), switchClaims.OrganizationID)

	// 切换到未归属组织 999 报错
	_, err = userSvc.SwitchOrganization(ctx, u.UserId, 999)
	require.Error(t, err)

	// 6. Profile 测试
	profile, err := userSvc.GetProfile(ctx, u.UserId)
	require.NoError(t, err)
	require.Equal(t, u.UserId, profile.UserId)

	require.NoError(t, userSvc.UpdateProfile(ctx, u.UserId, &v1.UpdateProfileRequest{Email: "u@e.com", Nickname: "neo"}))
}
