package internal

import (
	"time"
	"tradeengine/utils/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IDBBuilder interface {
	SetGormConfig(config *gorm.Config) IDBBuilder
	Build() (*gorm.DB, error)
}

type MySQLBuilder struct {
	connConfig   *mySQLConnConfig
	driverConfig *mysql.Config
	gormConfig   *gorm.Config
}

func NewMySQLBuilder(connConfig *mySQLConnConfig) *MySQLBuilder {
	return &MySQLBuilder{
		connConfig: connConfig,
	}
}

func (b *MySQLBuilder) SetDriverConfig(config *mysql.Config) {
	b.driverConfig = config
}

func (b *MySQLBuilder) getDriverConfig() *mysql.Config {
	if b.driverConfig == nil {
		b.driverConfig = b.getDefaultDriverConfig()
	}
	return b.driverConfig
}

func (b *MySQLBuilder) getDefaultDriverConfig() *mysql.Config {
	return &mysql.Config{
		DSN: b.getDSN(),
	}
}

func (b *MySQLBuilder) getDSN() string {
	return b.connConfig.getDSN()
}

func (b *MySQLBuilder) SetGormConfig(config *gorm.Config) IDBBuilder {
	b.gormConfig = config
	return b
}

func (b *MySQLBuilder) getGormConfig() *gorm.Config {
	if b.gormConfig == nil {
		b.gormConfig = b.getDefaultGormConfig()
	}
	return b.gormConfig
}

func (b *MySQLBuilder) getDefaultGormConfig() *gorm.Config {
	return &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (b *MySQLBuilder) getDialector() gorm.Dialector {
	return mysql.New(*b.getDriverConfig())
}

func (b *MySQLBuilder) Build() (*gorm.DB, error) {
	db, err := gorm.Open(b.getDialector(), b.getGormConfig())
	if err != nil {
		return nil, err
	}
	_, err = db.DB()
	if err != nil {
		logger.DB.Warn("failed to build MySQLBuilder, err: %s", err)
		return nil, err
	}

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	logger.DB.Warn("failed to build MySQLBuilder, err: %s", err)
	// 	return nil, err
	// }
	// sqlDB.SetMaxIdleConns(b.connConfig.MaxIdleConns)
	// sqlDB.SetMaxOpenConns(b.connConfig.MaxOpenConns)
	// sqlDB.SetConnMaxLifetime(b.connConfig.ConnMaxLifetime)
	return db, nil
}
