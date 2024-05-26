package db

import (
	"context"
	"sync"
	"tradeengine/service/db/internal"
)

var (
	DBMngr *DBManager
	once   sync.Once
)

type dbCategory uint64

type DBManager struct {
	ctx          context.Context
	dbMap        map[dbCategory]internal.DBService
	dbBuilderMap map[dbCategory]internal.IDBBuilder
}

func NewDBManager(ctx context.Context) *DBManager {
	once.Do(func() {
		DBMngr = &DBManager{
			ctx: ctx,
		}
	})
	return DBMngr
}

func (mngr *DBManager) Init() {

}

func (mngr *DBManager) Run() {
	go func() {

	}()
}
