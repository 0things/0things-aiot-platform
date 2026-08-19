package repository

import (
	"context"
	"fmt"
	"time"

	"aiot-backend/pkg/log"
	"aiot-backend/pkg/zapgorm2"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const ctxTxKey = "TxKey"

type Repository struct {
	db *gorm.DB
	//rdb    *redis.Client
	//mongo  *mongo.Client
	logger *log.Logger
}

func NewRepository(
	logger *log.Logger,
	db *gorm.DB,
	// rdb *redis.Client,
	//
	//	mongo *mongo.Client,
) *Repository {
	return &Repository{
		db: db,
		//rdb:    rdb,
		//mongo:  mongo,
		logger: logger,
	}
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTransaction(r *Repository) Transaction {
	return r
}

// DB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxTxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxTxKey, tx)
		return fn(ctx)
	})
}

func NewDB(conf *viper.Viper, l *log.Logger) *gorm.DB {
	return newDB(conf, l, "data.db.user")
}

// NewIoTDB opens the independent legacy device-service database. Keeping it
// separate prevents this migration from changing the user service datasource.
type IoTDB struct{ *gorm.DB }

func NewIoTDB(conf *viper.Viper, l *log.Logger) *IoTDB {
	return &IoTDB{newDB(conf, l, "data.db.device")}
}

func newDB(conf *viper.Viper, l *log.Logger, key string) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	logger := zapgorm2.New(l.Logger)
	driver := conf.GetString(key + ".driver")
	dsn := conf.GetString(key + ".dsn")
	if driver == "" || dsn == "" {
		panic("database configuration is missing for " + key)
	}

	// GORM doc: https://gorm.io/docs/connecting_to_the_database.html
	switch driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger,
		})
	case "postgres":
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		panic("unknown db driver")
	}
	if err != nil {
		panic(err)
	}
	db = db.Debug()

	// Connection Pool config
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	return rdb
}

type IoTRedis struct{ *redis.Client }

// NewIoTRedis is intentionally lazy: device CRUD must remain usable when
// telemetry storage is unavailable, just as the legacy service did.
func NewIoTRedis(conf *viper.Viper) *IoTRedis {
	addr := conf.GetString("data.redis.device.addr")
	if addr == "" {
		return &IoTRedis{}
	}
	return &IoTRedis{redis.NewClient(&redis.Options{
		Addr: addr, Password: conf.GetString("data.redis.device.password"), DB: conf.GetInt("data.redis.device.db"),
	})}
}
func NewMongo(conf *viper.Viper) (*mongo.Client, func(), error) {
	// https://www.mongodb.com/zh-cn/docs/drivers/go/current/
	uri := conf.GetString("data.mongo.uri")
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(fmt.Sprintf("mongo client error: %s", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Sprintf("mongo ping error: %s", err.Error()))
	}

	return client, func() {
		err = client.Disconnect(ctx)
		if err != nil {
			panic(fmt.Sprintf("mongo disconnect error: %s", err.Error()))
		}
	}, err
}
