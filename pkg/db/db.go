package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"jwt-authentication-golang/configs"
)

type DB struct {
	Con  *gorm.DB
	Info Info
}

type Info struct {
	DBName    string
	TableName string
}

var ConPool = &DB{}

func NewDBConPool(ctx context.Context, configs *configs.Configs) error {
	cfgs := configs
	ConPool.Info = Info{
		DBName:    cfgs.Database.Name,
		TableName: cfgs.Database.TableName,
	}
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		cfgs.Database.UserName,
		cfgs.Database.UserPassword,
		cfgs.Database.HostName,
		cfgs.Database.Port,
		cfgs.Database.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	ConPool.Con = db.WithContext(ctx)
	mysqlDB, err := db.DB()
	if err != nil {
		return err
	}
	err = mysqlDB.Ping()
	if err != nil {
		panic(err)
	}
	mysqlDB.SetConnMaxIdleTime(24 * time.Hour)
	mysqlDB.SetMaxOpenConns(cfgs.Database.MaxOpenCon)
	mysqlDB.SetMaxIdleConns(cfgs.Database.MaxIdleCon)
	return nil
}
