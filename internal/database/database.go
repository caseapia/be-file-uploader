package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/slog"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type Database struct {
	Web   *bun.DB
	Redis *redis.Client
}

func connectRedis() (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	password := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           0,
		PoolSize:     10,
		MaxIdleConns: 5,
		MaxRetries:   3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	slog.Infof("Successfully connected to redis: [%s]", addr)
	return rdb, nil
}

func connect(dbName string, maxOpen, maxIdle int) (*bun.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, password, host, port, dbName)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db %s: %w", dbName, err)
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(time.Minute * 3)
	sqlDB.SetConnMaxIdleTime(time.Minute * 1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("failed to ping db %s: %w", dbName, err)
	}

	db := bun.NewDB(sqlDB, mysqldialect.New())
	slog.Infof("Successfully connected to database: [%s] (MaxOpen: %d, MaxIdle: %d)", dbName, maxOpen, maxIdle)
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(os.Getenv("APP_MODE") == "DEV"),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return db, nil
}

func CreateDatabase() (*Database, error) {
	webDB := os.Getenv("DB_NAME")

	web, err := connect(webDB, 25, 10)
	if err != nil {
		return nil, err
	}

	rdb, err := connectRedis()
	if err != nil {
		return nil, err
	}

	return &Database{
		Web:   web,
		Redis: rdb,
	}, nil
}

func (d Database) Close() {
	if d.Web != nil {
		if err := d.Web.Close(); err != nil {
			slog.Errorf("failed to close database connection: %v", err)
		}
	}

	if d.Redis != nil {
		if err := d.Redis.Close(); err != nil {
			slog.Errorf("failed to close redis connection: %v", err)
		}
	}
}
