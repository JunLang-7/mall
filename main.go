package main

import (
	"context"
	"errors"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/config"
	"github.com/JunLang-7/mall/router"
	"github.com/JunLang-7/mall/service/user"
	"github.com/JunLang-7/mall/utils/logger"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	conf := config.InitConfig()
	logger.SetLevel(conf.Server.LogLevel)

	dbClient, err := initMySQL(&conf.MySQL)
	handleErr(err)
	logger.Debug("mysql connect success")

	rdsClient, err := initRedis(&conf.Redis)
	handleErr(err)
	logger.Debug("redis connect success")

	startServer(conf, dbClient, rdsClient).Run()
}

func startServer(conf *config.Config, db *gorm.DB, rds *redis.Client) *router.App {
	adpt := adaptor.NewAdaptor(conf, db, rds)
	user.NewService(adpt).StartBackgroundTasks(context.Background())
	return router.NewApp(
		conf.Server.HttpPort,
		router.NewRouter(
			conf,
			adpt,
			func() error {
				err := func() error {
					pingDb, err := db.DB()
					handleErr(err)
					return pingDb.Ping()
				}()
				if err != nil {
					return errors.New("MySQL connect failed")
				}
				return rds.Ping().Err()
			},
		),
	)
}

func initMySQL(conf *config.MySQL) (*gorm.DB, error) {
	dsn := conf.GetDsn()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(conf.MaxIdle)
	sqlDB.SetMaxOpenConns(conf.MaxOpen)
	return db, err
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func initRedis(conf *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.PWD,
		DB:           conf.DBIndex,
		MinIdleConns: conf.MaxIdle,
		PoolSize:     conf.MaxOpen,
	})
	if r, _ := client.Ping().Result(); r != "PONG" {
		return nil, errors.New("redis connect failed")
	}
	return client, nil
}
