package internal

import (
	"time"
	"tradeengine/service/db/model"
	"tradeengine/utils/logger"
	"tradeengine/utils/panichandle"
)

type DBName string
type DBMap map[DBName]*DBService

type IDBFactory interface {
	GetRequiredDBServiceMap() DBMap
	GetOptionalDBServiceMap() DBMap
}

func GetDBFactory() IDBFactory {
	if true {
		defaultDBFactory := NewDefaulDBFactory()
		return defaultDBFactory
	}
	return nil
}

type DefaulDBFactory struct {
	requiredDBServiceMap DBMap
	optionalDBServiceMap DBMap
}

func NewDefaulDBFactory() *DefaulDBFactory {
	factory := &DefaulDBFactory{
		requiredDBServiceMap: make(DBMap),
		optionalDBServiceMap: make(DBMap),
	}
	factory.initialize()
	return factory
}

func (f *DefaulDBFactory) GetRequiredDBServiceMap() DBMap {
	return f.requiredDBServiceMap
}

func (f *DefaulDBFactory) GetOptionalDBServiceMap() DBMap {
	return f.optionalDBServiceMap
}

func (f *DefaulDBFactory) initialize() {
	f.loadRequiredDBServiceMap()
	f.loadOptionalDBServiceMap()
}

func (f *DefaulDBFactory) loadRequiredDBServiceMap() {
	defer panichandle.PanicHandle()
	// FIXME: change to load information from env and files
	username := "trade_engine_admin"
	password := "trade_engine_is_666"
	var dbName DBName = "masterDB"
	connConfig := NewMySQLConnConfig("", "", username, password, string(dbName))
	connConfig.
		SetCharset("utf8mb4").
		SetParseTime("True").
		SetLoc("Local").
		SetTimeout("20s").
		SetReadTimeout("60s").
		SetWriteTimeout("60s")
	builder := NewMySQLDBBuilder(connConfig)
	modelList := []interface{}{model.Member{}}
	dbServ := NewDBService(builder, modelList)
	dbServ.SetMaxIdleConns(20)
	dbServ.SetMaxOpenConns(200)
	dbServ.SetConnMaxLifetime(time.Hour)
	if err := dbServ.Initialize(); err != nil {
		logger.DB.Error("failed to initial db service, err: %s", err)
		panic(err)
	}
	f.requiredDBServiceMap[dbName] = dbServ
}

func (f *DefaulDBFactory) loadOptionalDBServiceMap() {
	// TODO: reserved to implemnt optional DB services
}
