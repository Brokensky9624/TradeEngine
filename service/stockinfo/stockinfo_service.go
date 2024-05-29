package stockinfo

import (
	"sync"

	"tradeengine/server/web/rest/param"
	"tradeengine/service/db/model"
	dbTypes "tradeengine/service/db/types"
	"tradeengine/service/stockinfo/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"tradeengine/utils/logger"
	"tradeengine/utils/tool"
)

var (
	stockInfoSrv *StockInfoService
	once         sync.Once
)

func NewService(db *dbTypes.DBService) *StockInfoService {
	once.Do(func() {
		stockInfoSrv = &StockInfoService{
			db: db,
		}
	})
	return stockInfoSrv
}

func GetService() *StockInfoService {
	return stockInfoSrv
}

type StockInfoService struct {
	db *dbTypes.DBService
}

// public
func (s *StockInfoService) Create(param param.StockInfoCreateParam) error {
	var errPreFix string = "failed to create stock info"

	// check step
	if err := param.Check(); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}

	// create stock info
	createModel := &model.StockInfo{
		Name: param.Name,
	}
	if err := s.db.Create(createModel).Error; err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("stock info %d has created!", createModel.ID)
	return nil
}

func (s *StockInfoService) StockInfo(param param.StockInfoParam) (*types.StockInfo, error) {
	var errPreFix string = "failed to get stock info"

	// get stock info
	findModel := &model.StockInfo{}
	findModel.ID = param.ID
	_, err := s.takeStockInfoByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	logger.SERVER.Info("stock info %d has got!", findModel.ID)
	stockInfo := StockInfoModelToStockInfo(*findModel)
	return &stockInfo, nil
}

func (s *StockInfoService) StockInfoList(param param.StockInfoListParam) ([]types.StockInfo, error) {
	var errPreFix string = "failed to get stock info list"

	// find stock info list
	findModel := &model.StockInfo{}
	var modelInfoList []model.StockInfo
	tx := s.db.Where(*findModel).
		Order(clause.OrderByColumn{Column: clause.Column{Name: param.OrderBy}, Desc: param.OrderDesc}).
		Find(&modelInfoList)
	if tx.Error != nil {
		err := tool.PrefixError(errPreFix, tx.Error)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	stockInfoList := StockInfoModelListToStockInfoList(modelInfoList)
	logger.SERVER.Info("stock info list has got, len: %d !", len(stockInfoList))
	return stockInfoList, nil
}

// public: buy

func StockInfoModelToStockInfo(m model.StockInfo) types.StockInfo {
	stockInfo := types.StockInfo{
		ID:   m.ID,
		Name: m.Name,
	}
	return stockInfo
}

func StockInfoModelListToStockInfoList(modelList []model.StockInfo) []types.StockInfo {
	stockInfoList := make([]types.StockInfo, len(modelList))
	for i, m := range modelList {
		stockInfo := types.StockInfo{
			ID:   m.ID,
			Name: m.Name,
		}
		stockInfoList[i] = stockInfo
	}
	return stockInfoList
}

// private

func (s *StockInfoService) takeStockInfoByModel(findModel *model.StockInfo) (*gorm.DB, error) {
	tx := s.db.Where(*findModel).Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *StockInfoService) takeStockInfoBySeqAndModel(seq *gorm.DB, findModel *model.StockInfo) (*gorm.DB, error) {
	tx := seq.Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}
