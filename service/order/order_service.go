package order

import (
	"fmt"
	"sync"

	"tradeengine/server/web/rest/param"
	"tradeengine/service/db/model"
	dbTypes "tradeengine/service/db/types"
	"tradeengine/service/order/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"tradeengine/utils/logger"
	"tradeengine/utils/tool"
)

var (
	orderSrv *OrderService
	once     sync.Once
)

func NewService(db *dbTypes.DBService) *OrderService {
	once.Do(func() {
		orderSrv = &OrderService{
			db: db,
		}
	})
	return orderSrv
}

func GetService() *OrderService {
	return orderSrv
}

type OrderService struct {
	db *dbTypes.DBService
}

// public

func (s *OrderService) Create(param param.OrderCreateParam) error {
	var errPreFix string = "failed to create order"

	// check step
	if err := param.Check(); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}

	// create order
	orderType := types.OrderType(param.OrderType)
	if orderType == types.OrderTypeBuy {
		if err := s.createBuyOrder(param); err != nil {
			return err
		}
	} else if orderType == types.OrderTypeSell {
		if err := s.createSellOrder(param); err != nil {
			return err
		}
	} else {
		err := fmt.Errorf("%s, err: order's type is not supported", errPreFix)
		logger.SERVER.Debug(err.Error())
		return err
	}

	return nil
}

func (s *OrderService) Edit(param param.OrderEditParam) error {
	var errPreFix string = "failed to edit order"

	// check step
	if err := param.Check(); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}

	// update order
	orderType := types.OrderType(param.OrderType)
	if orderType == types.OrderTypeBuy {
		if err := s.editBuyOrder(param); err != nil {
			return err
		}
	} else if orderType == types.OrderTypeSell {
		if err := s.editSellOrder(param); err != nil {
			return err
		}
	} else {
		err := fmt.Errorf("%s, err: order's type is not supported", errPreFix)
		logger.SERVER.Debug(err.Error())
		return err
	}
	return nil
}

func (s *OrderService) Delete(param param.OrderDeleteParam) error {
	var errPreFix string = "failed to delete order"

	// check step
	if err := param.Check(); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}

	// delete order
	orderType := types.OrderType(param.OrderType)
	if orderType == types.OrderTypeBuy {
		if err := s.deleteBuyOrder(param); err != nil {
			return err
		}
	} else if orderType == types.OrderTypeSell {
		if err := s.deleteSellOrder(param); err != nil {
			return err
		}
	} else {
		err := fmt.Errorf("%s, err: order's type is not supported", errPreFix)
		logger.SERVER.Debug(err.Error())
		return err
	}
	return nil
}

func (s *OrderService) OrderInfo(param param.OrderInfoParam) (*types.Order, error) {
	var errPreFix string = "failed to get order info"

	// get order info
	orderType := types.OrderType(param.OrderType)
	if orderType == types.OrderTypeBuy {
		order, err := s.getBuyOrderInfo(param)
		if err != nil {
			return nil, err
		}
		return order, nil
	} else if orderType == types.OrderTypeSell {
		order, err := s.getSellOrderInfo(param)
		if err != nil {
			return nil, err
		}
		return order, nil
	} else {
		err := fmt.Errorf("%s, err: order's type is not supported", errPreFix)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
}

func (s *OrderService) OrderInfoList(param param.OrderInfoListParam) ([]types.Order, error) {
	var errPreFix string = "failed to get order list"

	// find order list
	orderType := types.OrderType(param.OrderType)
	if orderType == types.OrderTypeBuy {
		orderList, err := s.getBuyOrderInfoList(param)
		if err != nil {
			return nil, err
		}
		return orderList, nil
	} else if orderType == types.OrderTypeSell {
		orderList, err := s.getSellOrderInfoList(param)
		if err != nil {
			return nil, err
		}
		return orderList, nil
	} else {
		err := fmt.Errorf("%s, err: order's type is not supported", errPreFix)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
}

// public: buy

func BuyOrderModelToOrder(m model.BuyOrder) types.Order {
	order := types.Order{
		ID:          m.ID,
		OrderType:   uint(types.OrderTypeBuy),
		Price:       m.Price,
		OwnerID:     m.OwnerID,
		Quantity:    m.Quantity,
		StockInfoID: m.StockInfoID,
	}
	return order
}

func BuyOrderModelListToOrderList(modelList []model.BuyOrder) []types.Order {
	orderInfoList := make([]types.Order, len(modelList))
	for i, m := range modelList {
		order := types.Order{
			ID:          m.ID,
			OrderType:   uint(types.OrderTypeBuy),
			Price:       m.Price,
			OwnerID:     m.OwnerID,
			Quantity:    m.Quantity,
			StockInfoID: m.StockInfoID,
		}
		orderInfoList[i] = order
	}
	return orderInfoList
}

// public: sell

func SellOrderModelToOrder(m model.SellOrder) types.Order {
	order := types.Order{
		ID:          m.ID,
		OrderType:   uint(types.OrderTypeSell),
		Price:       m.Price,
		OwnerID:     m.OwnerID,
		Quantity:    m.Quantity,
		StockInfoID: m.StockInfoID,
	}
	return order
}

func SellOrderModelListToOrderList(modelList []model.SellOrder) []types.Order {
	orderInfoList := make([]types.Order, len(modelList))
	for i, m := range modelList {
		order := types.Order{
			ID:          m.ID,
			OrderType:   uint(types.OrderTypeSell),
			Price:       m.Price,
			OwnerID:     m.OwnerID,
			Quantity:    m.Quantity,
			StockInfoID: m.StockInfoID,
		}
		orderInfoList[i] = order
	}
	return orderInfoList
}

// private: order buy

func (s *OrderService) createBuyOrder(param param.OrderCreateParam) error {
	var errPreFix = "failed to create buy order"
	createModel := &model.BuyOrder{
		Price:       param.Price,
		OwnerID:     param.OwnerID,
		Quantity:    param.Quantity,
		StockInfoID: param.StockInfoID,
	}
	if err := s.db.Create(createModel).Error; err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("buy order %d has created!", createModel.ID)
	return nil
}

func (s *OrderService) editBuyOrder(param param.OrderEditParam) error {
	var errPreFix = "failed to modify buy order"
	findModel := &model.BuyOrder{}
	findModel.ID = param.ID
	findModel.OwnerID = param.OwnerID
	updateModel := &model.BuyOrder{}
	updateModel.ID = param.ID
	updateModel.Price = param.Price
	updateModel.Quantity = param.Quantity
	if _, err := s.updateOneBuyOrderByModel(findModel, updateModel); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("buy order %d has modified!", findModel.ID)
	return nil
}

func (s *OrderService) deleteBuyOrder(param param.OrderDeleteParam) error {
	var errPreFix = "failed to delete buy order"
	findModel := &model.BuyOrder{}
	findModel.ID = param.ID
	findModel.OwnerID = param.OwnerID
	_, err := s.deleteOneBuyOrderByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("buy order %d has deleted!", findModel.ID)
	return nil
}

func (s *OrderService) getBuyOrderInfo(param param.OrderInfoParam) (*types.Order, error) {
	var errPreFix = "failed to get buy order info"
	findModel := &model.BuyOrder{}
	findModel.ID = param.ID
	_, err := s.takeBuyOrderByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	logger.SERVER.Info("buy order info %d has got!", findModel.ID)
	order := BuyOrderModelToOrder(*findModel)
	return &order, nil
}

func (s *OrderService) getBuyOrderInfoList(param param.OrderInfoListParam) ([]types.Order, error) {
	var errPreFix = "failed to get buy order info list"
	findModel := &model.BuyOrder{}
	findModel.OwnerID = param.OwnerID
	var modelInfoList []model.BuyOrder
	tx := s.db.Where(*findModel).
		Order(clause.OrderByColumn{Column: clause.Column{Name: param.OrderBy}, Desc: param.OrderDesc}).
		Find(&modelInfoList)
	if tx.Error != nil {
		err := tool.PrefixError(errPreFix, tx.Error)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	orderInfoList := BuyOrderModelListToOrderList(modelInfoList)
	logger.SERVER.Info("buy order info list has got, len: %d !", len(orderInfoList))
	return orderInfoList, nil
}

func (s *OrderService) takeBuyOrderByModel(findModel *model.BuyOrder) (*gorm.DB, error) {
	tx := s.db.Where(*findModel).Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) takeBuyOrderBySeqAndModel(seq *gorm.DB, findModel *model.BuyOrder) (*gorm.DB, error) {
	tx := seq.Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) updateOneBuyOrderByModel(findModel *model.BuyOrder, updateModel *model.BuyOrder) (*gorm.DB, error) {
	tx, err := s.takeBuyOrderByModel(findModel)
	if err != nil {
		return nil, err
	}
	tx, err = s.updateOneBuyOrderBySeqAndModel(tx, updateModel)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *OrderService) updateOneBuyOrderBySeqAndModel(seq *gorm.DB, updateModel *model.BuyOrder) (*gorm.DB, error) {
	tx := seq.Updates(updateModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) deleteOneBuyOrderByModel(findModel *model.BuyOrder) (*gorm.DB, error) {
	tx, err := s.takeBuyOrderByModel(findModel)
	if err != nil {
		return nil, err
	}
	tx, err = s.deleteOneBuyOrderBySeqAndModel(tx, findModel)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *OrderService) deleteOneBuyOrderBySeqAndModel(seq *gorm.DB, findModel *model.BuyOrder) (*gorm.DB, error) {
	tx := seq.Delete(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if _, err := s.deleteUnscopedOneBuyOrderByModel(findModel); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *OrderService) deleteUnscopedOneBuyOrderByModel(findModel *model.BuyOrder) (*gorm.DB, error) {
	tx := s.db.Unscoped().Delete(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

// private: order sell

func (s *OrderService) createSellOrder(param param.OrderCreateParam) error {
	var errPreFix = "failed to create sell order"
	createModel := &model.SellOrder{
		Price:       param.Price,
		OwnerID:     param.OwnerID,
		Quantity:    param.Quantity,
		StockInfoID: param.StockInfoID,
	}
	if err := s.db.Create(createModel).Error; err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("sell order %d has created!", createModel.ID)
	return nil
}

func (s *OrderService) editSellOrder(param param.OrderEditParam) error {
	var errPreFix = "failed to modify sell order"
	findModel := &model.SellOrder{}
	findModel.ID = param.ID
	findModel.OwnerID = param.OwnerID
	updateModel := &model.SellOrder{}
	updateModel.ID = param.ID
	updateModel.Price = param.Price
	updateModel.Quantity = param.Quantity
	_, err := s.updateOneSellOrderByModel(findModel, updateModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("sell order %d has modified!", findModel.ID)
	return nil
}

func (s *OrderService) deleteSellOrder(param param.OrderDeleteParam) error {
	var errPreFix = "failed to delete sell order"
	findModel := &model.SellOrder{}
	findModel.ID = param.ID
	findModel.OwnerID = param.OwnerID
	_, err := s.deleteOneSellOrderByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("sell order %d has deleted!", findModel.ID)
	return nil
}

func (s *OrderService) getSellOrderInfo(param param.OrderInfoParam) (*types.Order, error) {
	var errPreFix = "failed to get sell order info"
	findModel := &model.SellOrder{}
	findModel.ID = param.ID
	_, err := s.takeSellOrderByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	logger.SERVER.Info("sell order info %d has got!", findModel.ID)
	order := SellOrderModelToOrder(*findModel)
	return &order, nil
}

func (s *OrderService) takeSellOrderByModel(findModel *model.SellOrder) (*gorm.DB, error) {
	tx := s.db.Where(*findModel).Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) takeSellOrderBySeqAndModel(seq *gorm.DB, findModel *model.SellOrder) (*gorm.DB, error) {
	tx := seq.Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) updateOneSellOrderByModel(findModel *model.SellOrder, updateModel *model.SellOrder) (*gorm.DB, error) {
	tx, err := s.takeSellOrderByModel(findModel)
	if err != nil {
		return nil, err
	}
	if tx, err = s.updateOneSellOrderBySeqAndModel(tx, updateModel); err != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) updateOneSellOrderBySeqAndModel(seq *gorm.DB, updateModel *model.SellOrder) (*gorm.DB, error) {
	tx := seq.Updates(updateModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) deleteOneSellOrderByModel(findModel *model.SellOrder) (*gorm.DB, error) {
	tx, err := s.takeSellOrderByModel(findModel)
	if err != nil {
		return nil, err
	}
	tx, err = s.deleteOneSellOrderBySeqAndModel(tx, findModel)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *OrderService) deleteOneSellOrderBySeqAndModel(seq *gorm.DB, findModel *model.SellOrder) (*gorm.DB, error) {
	tx := seq.Delete(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if _, err := s.deleteUnscopedOneSellOrderByModel(findModel); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *OrderService) deleteUnscopedOneSellOrderByModel(findModel *model.SellOrder) (*gorm.DB, error) {
	tx := s.db.Unscoped().Delete(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) getSellOrderInfoList(param param.OrderInfoListParam) ([]types.Order, error) {
	var errPreFix = "failed to get sell order info list"
	findModel := &model.SellOrder{}
	findModel.OwnerID = param.OwnerID
	var modelInfoList []model.SellOrder
	tx := s.db.Where(*findModel).
		Order(clause.OrderByColumn{Column: clause.Column{Name: param.OrderBy}, Desc: param.OrderDesc}).
		Find(&modelInfoList)
	if tx.Error != nil {
		err := tool.PrefixError(errPreFix, tx.Error)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	orderInfoList := SellOrderModelListToOrderList(modelInfoList)
	logger.SERVER.Info("sell order info list has got, len: %d !", len(orderInfoList))
	return orderInfoList, nil
}
