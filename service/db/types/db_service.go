package types

import (
	"database/sql"
	"time"
	"tradeengine/service/db/internal"

	"gorm.io/gorm"
)

type DBService struct {
	*gorm.DB
	dbBuilder       internal.IDBBuilder
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	modelList       []interface{}
}

func NewDBService(dbBuilder internal.IDBBuilder, modelList []interface{}) *DBService {
	return &DBService{
		dbBuilder: dbBuilder,
		modelList: modelList,
	}
}

func (d *DBService) SetMaxIdleConns(maxIdleConns int) *DBService {
	d.MaxIdleConns = maxIdleConns
	return d
}

func (d *DBService) SetMaxOpenConns(maxOpenConns int) *DBService {
	d.MaxOpenConns = maxOpenConns
	return d
}

func (d *DBService) SetConnMaxLifetime(connMaxLifetime time.Duration) *DBService {
	d.ConnMaxLifetime = connMaxLifetime
	return d
}

func (d *DBService) Initialize() error {
	db, err := d.BuildDB()
	if err != nil {
		return err
	}
	d.DB = db
	if err = d.Ping(); err != nil {
		return err
	}
	if err = d.AutoMigrate(d.modelList...); err != nil {
		return err
	}
	return nil
}

func (d *DBService) BuildDB() (*gorm.DB, error) {
	db, err := d.dbBuilder.Build()
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(d.MaxIdleConns)
	sqlDB.SetMaxOpenConns(d.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(d.ConnMaxLifetime)
	return db, nil
}

func (d *DBService) Ping() error {
	sqlDB, err := d.GetSqlDB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (d *DBService) Close() error {
	sqlDB, err := d.GetSqlDB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *DBService) GetSqlDB() (*sql.DB, error) {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return nil, err
	}
	return sqlDB, nil
}
