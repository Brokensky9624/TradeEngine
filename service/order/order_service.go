package order

import (
	"fmt"
	"sync"

	"tradeengine/server/web/rest/param"
	model "tradeengine/service/db/model/order"
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
		if err := s.createOrderBuy(param); err != nil {
			return err
		}
	} else if orderType == types.OrderTypeSell {
		if err := s.createOrderSell(param); err != nil {
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
		if err := s.editOrderBuy(param); err != nil {
			return err
		}
	} else if orderType == types.OrderTypeSell {
		if err := s.editOrderSell(param); err != nil {
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
		if err := s.deleteOrderBuy(param); err != nil {
			return err
		}
	} else if orderType == types.OrderTypeSell {
		if err := s.deleteOrderSell(param); err != nil {
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
		order, err := s.getOrderBuyInfo(param)
		if err != nil {
			return nil, err
		}
		return order, nil
	} else if orderType == types.OrderTypeSell {
		order, err := s.getOrderSellInfo(param)
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
		orderList, err := s.getOrderBuyInfoList(param)
		if err != nil {
			return nil, err
		}
		return orderList, nil
	} else if orderType == types.OrderTypeSell {
		orderList, err := s.getOrderSellInfoList(param)
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

func OrderBuyModelToOrder(m model.OrderBuy) types.Order {
	order := types.Order{
		ID:        m.ID,
		OrderType: uint(types.OrderTypeBuy),
		Price:     m.Price,
		OwnerID:   m.OwnerID,
		Quantity:  m.Quantity,
		StockID:   m.StockID,
		Status:    m.Status,
	}
	return order
}

func OrderBuyModelListToOrderList(modelList []model.OrderBuy) []types.Order {
	orderInfoList := make([]types.Order, len(modelList))
	for i, m := range modelList {
		order := types.Order{
			ID:        m.ID,
			OrderType: uint(types.OrderTypeBuy),
			Price:     m.Price,
			OwnerID:   m.OwnerID,
			Quantity:  m.Quantity,
			StockID:   m.StockID,
			Status:    m.Status,
		}
		orderInfoList[i] = order
	}
	return orderInfoList
}

// public: sell

func OrderSellModelToOrder(m model.OrderSell) types.Order {
	order := types.Order{
		ID:        m.ID,
		OrderType: uint(types.OrderTypeSell),
		Price:     m.Price,
		OwnerID:   m.OwnerID,
		Quantity:  m.Quantity,
		StockID:   m.StockID,
		Status:    m.Status,
	}
	return order
}

func OrderSellModelListToOrderList(modelList []model.OrderSell) []types.Order {
	orderInfoList := make([]types.Order, len(modelList))
	for i, m := range modelList {
		order := types.Order{
			ID:        m.ID,
			OrderType: uint(types.OrderTypeSell),
			Price:     m.Price,
			OwnerID:   m.OwnerID,
			Quantity:  m.Quantity,
			StockID:   m.StockID,
			Status:    m.Status,
		}
		orderInfoList[i] = order
	}
	return orderInfoList
}

// private: order buy

func (s *OrderService) createOrderBuy(param param.OrderCreateParam) error {
	var errPreFix = "failed to create buy order"
	createModel := &model.OrderBuy{
		Price:    param.Price,
		OwnerID:  param.OwnerID,
		Quantity: param.Quantity,
		StockID:  param.StockID,
		Status:   uint(types.OrderStatusNew),
	}
	if err := s.db.Create(createModel).Error; err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("buy order %d has created!", createModel.ID)
	return nil
}

func (s *OrderService) editOrderBuy(param param.OrderEditParam) error {
	var errPreFix = "failed to modify buy order"
	findModel := &model.OrderBuy{}
	findModel.ID = param.ID
	findModel.OwnerID = param.OwnerID
	updateModel := &model.OrderBuy{}
	updateModel.ID = param.ID
	updateModel.Price = param.Price
	updateModel.Quantity = param.Quantity
	if _, err := s.updateOneOrderBuyByModel(findModel, updateModel); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("buy order %d has modified!", findModel.ID)
	return nil
}

func (s *OrderService) deleteOrderBuy(param param.OrderDeleteParam) error {
	var errPreFix = "failed to delete buy order"
	findModel := &model.OrderBuy{}
	findModel.ID = param.ID
	findModel.OwnerID = param.OwnerID
	_, err := s.deleteOneOrderBuyByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("buy order %d has deleted!", findModel.ID)
	return nil
}

func (s *OrderService) getOrderBuyInfo(param param.OrderInfoParam) (*types.Order, error) {
	var errPreFix = "failed to get buy order info"
	findModel := &model.OrderBuy{}
	findModel.ID = param.ID
	_, err := s.takeOrderBuyByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	logger.SERVER.Info("buy order info %d has got!", findModel.ID)
	order := OrderBuyModelToOrder(*findModel)
	return &order, nil
}

func (s *OrderService) getOrderBuyInfoList(param param.OrderInfoListParam) ([]types.Order, error) {
	var errPreFix = "failed to get buy order info list"
	findModel := &model.OrderBuy{}
	findModel.OwnerID = param.OwnerID
	var modelInfoList []model.OrderBuy
	tx := s.db.Where(*findModel).
		Order(clause.OrderByColumn{Column: clause.Column{Name: param.OrderBy}, Desc: param.OrderDesc}).
		Find(&modelInfoList)
	if tx.Error != nil {
		err := tool.PrefixError(errPreFix, tx.Error)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	orderInfoList := OrderBuyModelListToOrderList(modelInfoList)
	logger.SERVER.Info("buy order info list has got, len: %d !", len(orderInfoList))
	return orderInfoList, nil
}

func (s *OrderService) takeOrderBuyByModel(findModel *model.OrderBuy) (*gorm.DB, error) {
	tx := s.db.Where(*findModel).Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) takeOrderBuyBySeqAndModel(seq *gorm.DB, findModel *model.OrderBuy) (*gorm.DB, error) {
	tx := seq.Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) updateOneOrderBuyByModel(findModel *model.OrderBuy, updateModel *model.OrderBuy) (*gorm.DB, error) {
	tx, err := s.takeOrderBuyByModel(findModel)
	if err != nil {
		return nil, err
	}
	tx, err = s.updateOneOrderBuyBySeqAndModel(tx, updateModel)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *OrderService) updateOneOrderBuyBySeqAndModel(seq *gorm.DB, updateModel *model.OrderBuy) (*gorm.DB, error) {
	tx := seq.Updates(updateModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) deleteOneOrderBuyByModel(findModel *model.OrderBuy) (*gorm.DB, error) {
	tx, err := s.takeOrderBuyByModel(findModel)
	if err != nil {
		return nil, err
	}
	tx, err = s.deleteOneOrderBuyBySeqAndModel(tx, findModel)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *OrderService) deleteOneOrderBuyBySeqAndModel(seq *gorm.DB, findModel *model.OrderBuy) (*gorm.DB, error) {
	tx := seq.Delete(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if _, err := s.deleteUnscopedOneOrderBuyByModel(findModel); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *OrderService) deleteUnscopedOneOrderBuyByModel(findModel *model.OrderBuy) (*gorm.DB, error) {
	tx := s.db.Unscoped().Delete(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

// private: order sell

func (s *OrderService) createOrderSell(param param.OrderCreateParam) error {
	var errPreFix = "failed to create sell order"
	createModel := &model.OrderSell{
		Price:    param.Price,
		OwnerID:  param.OwnerID,
		Quantity: param.Quantity,
		StockID:  param.StockID,
		Status:   uint(types.OrderStatusNew),
	}
	if err := s.db.Create(createModel).Error; err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("sell order %d has created!", createModel.ID)
	return nil
}

func (s *OrderService) editOrderSell(param param.OrderEditParam) error {
	var errPreFix = "failed to modify sell order"
	findModel := &model.OrderSell{}
	findModel.ID = param.ID
	findModel.OwnerID = param.OwnerID
	updateModel := &model.OrderSell{}
	updateModel.ID = param.ID
	updateModel.Price = param.Price
	updateModel.Quantity = param.Quantity
	_, err := s.updateOneOrderSellByModel(findModel, updateModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("sell order %d has modified!", findModel.ID)
	return nil
}

func (s *OrderService) deleteOrderSell(param param.OrderDeleteParam) error {
	var errPreFix = "failed to delete sell order"
	findModel := &model.OrderSell{}
	findModel.ID = param.ID
	findModel.OwnerID = param.OwnerID
	_, err := s.deleteOneOrderSellByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("sell order %d has deleted!", findModel.ID)
	return nil
}

func (s *OrderService) getOrderSellInfo(param param.OrderInfoParam) (*types.Order, error) {
	var errPreFix = "failed to get sell order info"
	findModel := &model.OrderSell{}
	findModel.ID = param.ID
	_, err := s.takeOrderSellByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	logger.SERVER.Info("sell order info %d has got!", findModel.ID)
	order := OrderSellModelToOrder(*findModel)
	return &order, nil
}

func (s *OrderService) takeOrderSellByModel(findModel *model.OrderSell) (*gorm.DB, error) {
	tx := s.db.Where(*findModel).Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) takeOrderSellBySeqAndModel(seq *gorm.DB, findModel *model.OrderSell) (*gorm.DB, error) {
	tx := seq.Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) updateOneOrderSellByModel(findModel *model.OrderSell, updateModel *model.OrderSell) (*gorm.DB, error) {
	tx, err := s.takeOrderSellByModel(findModel)
	if err != nil {
		return nil, err
	}
	if tx, err = s.updateOneOrderSellBySeqAndModel(tx, updateModel); err != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) updateOneOrderSellBySeqAndModel(seq *gorm.DB, updateModel *model.OrderSell) (*gorm.DB, error) {
	tx := seq.Updates(updateModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) deleteOneOrderSellByModel(findModel *model.OrderSell) (*gorm.DB, error) {
	tx, err := s.takeOrderSellByModel(findModel)
	if err != nil {
		return nil, err
	}
	tx, err = s.deleteOneOrderSellBySeqAndModel(tx, findModel)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *OrderService) deleteOneOrderSellBySeqAndModel(seq *gorm.DB, findModel *model.OrderSell) (*gorm.DB, error) {
	tx := seq.Delete(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if _, err := s.deleteUnscopedOneOrderSellByModel(findModel); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *OrderService) deleteUnscopedOneOrderSellByModel(findModel *model.OrderSell) (*gorm.DB, error) {
	tx := s.db.Unscoped().Delete(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *OrderService) getOrderSellInfoList(param param.OrderInfoListParam) ([]types.Order, error) {
	var errPreFix = "failed to get sell order info list"
	findModel := &model.OrderSell{}
	findModel.OwnerID = param.OwnerID
	var modelInfoList []model.OrderSell
	tx := s.db.Where(*findModel).
		Order(clause.OrderByColumn{Column: clause.Column{Name: param.OrderBy}, Desc: param.OrderDesc}).
		Find(&modelInfoList)
	if tx.Error != nil {
		err := tool.PrefixError(errPreFix, tx.Error)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	orderInfoList := OrderSellModelListToOrderList(modelInfoList)
	logger.SERVER.Info("sell order info list has got, len: %d !", len(orderInfoList))
	return orderInfoList, nil
}
