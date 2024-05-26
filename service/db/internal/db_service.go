package internal

import "gorm.io/gorm"

type DBService struct {
	*gorm.DB
}

func NewDBService(db *gorm.DB) *DBService {
	return &DBService{
		DB: db,
	}
}

func (d *DBService) InitTable() error {
	return d.AutoMigrate(
	// &model.Member{},
	// &model.Product{},
	)
}

func (d *DBService) Ping() error {
	// get sqlDB for close later
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (d *DBService) Close() error {
	// get sqlDB for close later
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
