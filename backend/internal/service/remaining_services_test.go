package service

import (
	"context"
	"testing"

	v1 "aiot-backend/api/v1"
	messageParserV1 "aiot-backend/api/message_parser/v1"
	"aiot-backend/internal/model"
	"aiot-backend/internal/repository"

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
	require.NoError(t, db.AutoMigrate(&model.User{}))

	logger := log.NewLog(viper.New())
	repo := repository.NewRepository(logger, db)
	tm := repository.NewTransaction(repo)
	svc := NewService(tm, logger, sid.NewSid(), jwt.NewJwt(viper.New()))
	userSvc := NewUserService(svc, repository.NewUserRepository(repo))
	ctx := context.Background()

	require.NoError(t, userSvc.Register(ctx, &v1.RegisterRequest{Email: "u@e.com", Password: "secret"}))
	require.Error(t, userSvc.Register(ctx, &v1.RegisterRequest{Email: "u@e.com", Password: "secret"}))

	token, err := userSvc.Login(ctx, &v1.LoginRequest{Email: "u@e.com", Password: "secret"})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// wrong password
	_, err = userSvc.Login(ctx, &v1.LoginRequest{Email: "u@e.com", Password: "wrong"})
	require.Error(t, err)

	var u model.User
	require.NoError(t, db.Where("email = ?", "u@e.com").First(&u).Error)
	require.NotEmpty(t, u.UserId)

	profile, err := userSvc.GetProfile(ctx, u.UserId)
	require.NoError(t, err)
	require.Equal(t, u.UserId, profile.UserId)

	require.NoError(t, userSvc.UpdateProfile(ctx, u.UserId, &v1.UpdateProfileRequest{Email: "u@e.com", Nickname: "neo"}))
}
