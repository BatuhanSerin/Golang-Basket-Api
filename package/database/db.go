package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/BatuhanSerin/final-project/package/config"
	"go.uber.org/zap"
	gormPsql "gorm.io/driver/postgres"
)

func Connect(cfg *config.Config) *gorm.DB {

	db, err := gorm.Open(gormPsql.Open(cfg.DatabaseConfig.DataSourceName), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("DB Connection Error", zap.Error(err))
	}
	origin, err := db.DB()
	if err != nil {
		zap.L().Fatal("DB Connection Error", zap.Error(err))
	}
	// err = origin.Ping()
	// if err != nil {
	// 	zap.L().Fatal("DB Connection Error", zap.Error(err))
	// }

	origin.SetMaxIdleConns(cfg.DatabaseConfig.MaxIdle)
	origin.SetMaxOpenConns(cfg.DatabaseConfig.MaxOpen)
	origin.SetConnMaxLifetime(time.Duration(cfg.DatabaseConfig.MaxLifetime) * time.Second)

	return db
}
