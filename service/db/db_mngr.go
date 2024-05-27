package db

import (
	"context"
	"sync"
	"tradeengine/service/db/factory"
	"tradeengine/service/db/types"
	"tradeengine/utils/logger"
)

const defaultDB = "tradeEngineDB"

var (
	DBMngr *DBManager
	once   sync.Once
)

type DBManager struct {
	ctx                  context.Context
	factory              factory.IDBFactory
	requiredDBServiceMap factory.DBMap
	optionalDBServiceMap factory.DBMap
}

func NewDBManager(ctx context.Context) *DBManager {
	once.Do(func() {
		DBMngr = &DBManager{
			ctx: ctx,
		}
		DBMngr.initialize()
	})
	return DBMngr
}

func (mngr *DBManager) initialize() {
	logger.DB.Info("Initializing DBManager")
	mngr.factory = factory.GetDBFactory()
	mngr.loadDBServiceMap()
	logger.DB.Info("DBManager is ready")

}

func (mngr *DBManager) loadDBServiceMap() {
	mngr.requiredDBServiceMap = mngr.factory.GetRequiredDBServiceMap()
	mngr.optionalDBServiceMap = mngr.factory.GetOptionalDBServiceMap()
}

func (mngr *DBManager) DefaultDBService() *types.DBService {
	return mngr.DBService(defaultDB)
}

func (mngr *DBManager) DBService(dbName string) *types.DBService {
	formatedDBName := factory.DBName(dbName)
	if v, ok := mngr.requiredDBServiceMap[formatedDBName]; ok {
		return v
	} else if v, ok := mngr.optionalDBServiceMap[formatedDBName]; ok {
		return v
	}
	return nil
}

func GetMngr() *DBManager {
	return DBMngr
}

func (mngr *DBManager) Run() {
	go func() {
		<-mngr.ctx.Done()
		for _, requiredDbSrv := range mngr.requiredDBServiceMap {
			requiredDbSrv.Close()
		}
		for _, optionalDBSrv := range mngr.optionalDBServiceMap {
			optionalDBSrv.Close()
		}
		logger.SERVER.Info("all db were closed")
	}()
}
