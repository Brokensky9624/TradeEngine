package factory

import (
	"time"
	"tradeengine/service/db/internal"
	"tradeengine/service/db/model"
	"tradeengine/service/db/types"
	"tradeengine/utils/logger"
	"tradeengine/utils/panichandle"
)

type DBName string
type DBMap map[DBName]*types.DBService

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
	var dbName DBName = "tradeEngineDB"
	connConfig := internal.NewMySQLConnConfig("", "", username, password, string(dbName))
	connConfig.
		SetCharset("utf8mb4").
		SetParseTime("True").
		SetLoc("Local").
		SetTimeout("20s").
		SetReadTimeout("60s").
		SetWriteTimeout("60s")
	builder := internal.NewMySQLDBBuilder(connConfig)
	modelList := []interface{}{
		model.Member{},
		model.Wallet{},
		model.StockInfo{},
		model.Stock{},
		model.BuyOrder{},
		model.SellOrder{},
	}
	dbServ := types.NewDBService(builder, modelList)
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
