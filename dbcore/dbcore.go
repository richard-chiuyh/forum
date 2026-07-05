package dbcore

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type MysqlConfig struct {
	DataSource   []string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifeTime  int
}

func Connect(conf *MysqlConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.DataSource[0]), &gorm.Config{Logger: &Logger{}})
	if err != nil {
		return nil, err
	}
	// 添加链路追踪
	if err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		return nil, err
	}
	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqldb.SetMaxIdleConns(conf.MaxIdleConns)
	sqldb.SetMaxOpenConns(conf.MaxOpenConns)
	sqldb.SetConnMaxLifetime(time.Duration(conf.MaxLifeTime) * time.Second)
	return db, nil
}
