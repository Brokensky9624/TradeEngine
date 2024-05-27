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

type MySQLDBBuilder struct {
	connConfig   *mySQLConnConfig
	driverConfig *mysql.Config
	gormConfig   *gorm.Config
}

func NewMySQLDBBuilder(connConfig *mySQLConnConfig) *MySQLDBBuilder {
	return &MySQLDBBuilder{
		connConfig: connConfig,
	}
}

func (b *MySQLDBBuilder) SetDriverConfig(config *mysql.Config) {
	b.driverConfig = config
}

func (b *MySQLDBBuilder) getDriverConfig() *mysql.Config {
	if b.driverConfig == nil {
		b.driverConfig = b.getDefaultDriverConfig()
	}
	return b.driverConfig
}

func (b *MySQLDBBuilder) getDefaultDriverConfig() *mysql.Config {
	return &mysql.Config{
		DSN: b.getDSN(),
	}
}

func (b *MySQLDBBuilder) getDSN() string {
	return b.connConfig.getDSN()
}

func (b *MySQLDBBuilder) SetGormConfig(config *gorm.Config) IDBBuilder {
	b.gormConfig = config
	return b
}

func (b *MySQLDBBuilder) getGormConfig() *gorm.Config {
	if b.gormConfig == nil {
		b.gormConfig = b.getDefaultGormConfig()
	}
	return b.gormConfig
}

func (b *MySQLDBBuilder) getDefaultGormConfig() *gorm.Config {
	return &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (b *MySQLDBBuilder) getDialector() gorm.Dialector {
	return mysql.New(*b.getDriverConfig())
}

func (b *MySQLDBBuilder) Build() (*gorm.DB, error) {
	db, err := gorm.Open(b.getDialector(), b.getGormConfig())
	if err != nil {
		return nil, err
	}
	_, err = db.DB()
	if err != nil {
		logger.DB.Error("failed to build MySQLDBBuilder, err: %s", err)
		return nil, err
	}
	return db, nil
}
