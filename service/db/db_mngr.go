package db

import (
	"context"
	"sync"
	"tradeengine/service/db/internal"
	"tradeengine/utils/logger"
)

var (
	DBMngr *DBManager
	once   sync.Once
)

type DBManager struct {
	ctx                  context.Context
	factory              internal.IDBFactory
	requiredDBServiceMap internal.DBMap
	optionalDBServiceMap internal.DBMap
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
	mngr.factory = internal.GetDBFactory()
	mngr.loadDBServiceMap()
	logger.DB.Info("DBManager is ready")

}

func (mngr *DBManager) loadDBServiceMap() {
	mngr.requiredDBServiceMap = mngr.factory.GetRequiredDBServiceMap()
	mngr.optionalDBServiceMap = mngr.factory.GetOptionalDBServiceMap()
}

func (mngr *DBManager) DBService(dbName string) *internal.DBService {
	formatedDBName := internal.DBName(dbName)
	if v, ok := mngr.requiredDBServiceMap[formatedDBName]; ok {
		return v
	} else if v, ok := mngr.optionalDBServiceMap[formatedDBName]; ok {
		return v
	}
	return nil
}

func (mngr *DBManager) Run() <-chan interface{} {
	finishedChan := make(chan interface{})
	go func() {
		defer close(finishedChan)
		for {
			select {
			case <-mngr.ctx.Done():
				for _, requiredDbSrv := range mngr.requiredDBServiceMap {
					requiredDbSrv.Close()
				}
				for _, optionalDBSrv := range mngr.optionalDBServiceMap {
					optionalDBSrv.Close()
				}
				return
			default:
				continue
			}
		}
	}()
	return finishedChan
}
