package stock

import (
	"sync"

	"tradeengine/server/web/rest/param"
	"tradeengine/service/db/model"
	dbTypes "tradeengine/service/db/types"
	"tradeengine/service/stock/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"tradeengine/utils/logger"
	"tradeengine/utils/tool"
)

var (
	stockSrv *StockService
	once     sync.Once
)

func NewService(db *dbTypes.DBService) *StockService {
	once.Do(func() {
		stockSrv = &StockService{
			db: db,
		}
	})
	return stockSrv
}

func GetService() *StockService {
	return stockSrv
}

type StockService struct {
	db *dbTypes.DBService
}

// public
func (s *StockService) Create(param param.OneStockCreateParam) error {
	var errPreFix string = "failed to create stock"

	// check step
	if err := param.Check(); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}

	// create stock
	createModel := &model.Stock{
		StockInfoID:       param.StockInfoID,
		OwnerID:           param.OwnerID,
		AvailableQuantity: param.AvailableQuantity,
		PendingQuantity:   param.PendingQuantity,
	}
	if err := s.db.Create(createModel).Error; err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("stock %d has created!", createModel.ID)
	return nil
}

func (s *StockService) Edit(param param.OneStockEditParam) error {
	var errPreFix = "failed to edit stock"
	findModel := &model.Stock{}
	findModel.ID = param.ID
	findModel.OwnerID = param.OwnerID
	updateModel := &model.Stock{}
	updateModel.ID = param.ID
	updateModel.AvailableQuantity = param.AvailableQuantity
	updateModel.PendingQuantity = param.PendingQuantity
	if _, err := s.updateOneStockByModel(findModel, updateModel); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("buy order %d has modified!", findModel.ID)
	return nil
}

func (s *StockService) Delete(param param.OneStockDeleteParam) error {
	var errPreFix = "failed to delete buy order"
	findModel := &model.Stock{}
	findModel.ID = param.ID
	findModel.OwnerID = param.OwnerID
	_, err := s.deleteOneStockByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("delete stock %d has deleted!", findModel.ID)
	return nil
}

func (s *StockService) OneStockInfo(param param.OneStockInfoParam) (*types.Stock, error) {
	var errPreFix string = "failed to get stock"

	// get stock
	findModel := &model.Stock{}
	findModel.ID = param.ID
	_, err := s.takeStockByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	logger.SERVER.Info("stock %d has got!", findModel.ID)
	stock := StockModelToStock(*findModel)
	return &stock, nil
}

func (s *StockService) OneStockInfoList(param param.OneStockInfoListParam) ([]types.Stock, error) {
	var errPreFix string = "failed to get stock list"

	// find stock list
	findModel := &model.Stock{}
	var modelInfoList []model.Stock
	tx := s.db.Where(*findModel).
		Order(clause.OrderByColumn{Column: clause.Column{Name: param.OrderBy}, Desc: param.OrderDesc}).
		Find(&modelInfoList)
	if tx.Error != nil {
		err := tool.PrefixError(errPreFix, tx.Error)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	stockList := StockModelListToStockList(modelInfoList)
	logger.SERVER.Info("stock list has got, len: %d !", len(stockList))
	return stockList, nil
}

// public: buy

func StockModelToStock(m model.Stock) types.Stock {
	stock := types.Stock{
		ID:                m.ID,
		OwnerID:           m.OwnerID,
		AvailableQuantity: m.AvailableQuantity,
		PendingQuantity:   m.PendingQuantity,
		StockInfoID:       m.StockInfoID,
	}
	return stock
}

func StockModelListToStockList(modelList []model.Stock) []types.Stock {
	stockList := make([]types.Stock, len(modelList))
	for i, m := range modelList {
		stock := StockModelToStock(m)
		stockList[i] = stock
	}
	return stockList
}

// private

func (s *StockService) takeStockByModel(findModel *model.Stock) (*gorm.DB, error) {
	tx := s.db.Where(*findModel).Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *StockService) takeStockBySeqAndModel(seq *gorm.DB, findModel *model.Stock) (*gorm.DB, error) {
	tx := seq.Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

// private

func (s *StockService) takeOneStockByModel(findModel *model.Stock) (*gorm.DB, error) {
	tx := s.db.Where(*findModel).Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *StockService) updateOneStockByModel(findModel *model.Stock, updateModel *model.Stock) (*gorm.DB, error) {
	tx, err := s.takeOneStockByModel(findModel)
	if err != nil {
		return nil, err
	}
	tx, err = s.updateOneStockBySeqAndModel(tx, updateModel)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *StockService) updateOneStockBySeqAndModel(seq *gorm.DB, updateModel *model.Stock) (*gorm.DB, error) {
	tx := seq.Updates(updateModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *StockService) deleteOneStockByModel(findModel *model.Stock) (*gorm.DB, error) {
	tx, err := s.takeOneStockByModel(findModel)
	if err != nil {
		return nil, err
	}
	tx, err = s.deleteOneStockBySeqAndModel(tx, findModel)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *StockService) deleteOneStockBySeqAndModel(seq *gorm.DB, findModel *model.Stock) (*gorm.DB, error) {
	tx := seq.Delete(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if _, err := s.deleteUnscopedOneStockByModel(findModel); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *StockService) deleteUnscopedOneStockByModel(findModel *model.Stock) (*gorm.DB, error) {
	tx := s.db.Unscoped().Delete(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}
