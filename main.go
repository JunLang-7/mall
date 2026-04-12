package main

import (
	"errors"

	"github.com/JunLang-7/mall/adaptor"
	"github.com/JunLang-7/mall/config"
	"github.com/JunLang-7/mall/router"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	conf := config.InitConfig()

	dbClient, err := initMySQL(&conf.MySQL)
	handleErr(err)

	startServer(conf, dbClient)
}

func startServer(conf *config.Config, db *gorm.DB) *router.App {
	return router.NewApp(conf.Server.HttpPort,
		router.NewRouter(
			conf,
			adaptor.NewAdaptor(conf, db),
			func() error {
				err := func() error {
					pingDb, err := db.DB()
					handleErr(err)
					return pingDb.Ping()
				}()
				if err != nil {
					return errors.New("MySQL connect failed")
				}
				return nil
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
	return db, err
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
