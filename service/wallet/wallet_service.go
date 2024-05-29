package wallet

import (
	"sync"

	"tradeengine/server/web/rest/param"
	"tradeengine/service/db/model"
	dbTypes "tradeengine/service/db/types"
	"tradeengine/service/wallet/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"tradeengine/utils/logger"
	"tradeengine/utils/tool"
)

var (
	walletSrv *WalletService
	once      sync.Once
)

func NewService(db *dbTypes.DBService) *WalletService {
	once.Do(func() {
		walletSrv = &WalletService{
			db: db,
		}
	})
	return walletSrv
}

func GetService() *WalletService {
	return walletSrv
}

type WalletService struct {
	db *dbTypes.DBService
}

// public
func (s *WalletService) Create(param param.WalletCreateParam) error {
	var errPreFix string = "failed to create wallet"

	// check step
	if err := param.Check(); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}

	// create wallet
	createModel := &model.Wallet{
		OwnerID:        param.OwnerID,
		AvailableMoney: param.AvailableMoney,
		PendingMoney:   param.PendingMoney,
	}
	if err := s.db.Create(createModel).Error; err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("wallet %d has created!", createModel.ID)
	return nil
}

func (s *WalletService) Edit(param param.WalletEditParam) error {
	var errPreFix = "failed to edit wallet"
	findModel := &model.Wallet{}
	findModel.ID = param.ID
	findModel.OwnerID = param.OwnerID
	updateModel := &model.Wallet{}
	updateModel.ID = param.ID
	updateModel.AvailableMoney = param.AvailableMoney
	updateModel.PendingMoney = param.PendingMoney
	if _, err := s.updateOneWalletByModel(findModel, updateModel); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("buy order %d has modified!", findModel.ID)
	return nil
}

func (s *WalletService) WalletInfo(param param.WalletInfoParam) (*types.Wallet, error) {
	var errPreFix string = "failed to get wallet"

	// get wallet
	findModel := &model.Wallet{}
	findModel.ID = param.ID
	_, err := s.takeWalletByModel(findModel)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	logger.SERVER.Info("wallet %d has got!", findModel.ID)
	wallet := WalletModelToWallet(*findModel)
	return &wallet, nil
}

func (s *WalletService) WalletInfoList(param param.WalletInfoListParam) ([]types.Wallet, error) {
	var errPreFix string = "failed to get wallet list"

	// find wallet list
	findModel := &model.Wallet{}
	var modelInfoList []model.Wallet
	tx := s.db.Where(*findModel).
		Order(clause.OrderByColumn{Column: clause.Column{Name: param.OrderBy}, Desc: param.OrderDesc}).
		Find(&modelInfoList)
	if tx.Error != nil {
		err := tool.PrefixError(errPreFix, tx.Error)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	walletList := WalletModelListToWalletList(modelInfoList)
	logger.SERVER.Info("wallet list has got, len: %d !", len(walletList))
	return walletList, nil
}

// public: buy

func WalletModelToWallet(m model.Wallet) types.Wallet {
	wallet := types.Wallet{
		ID:             m.ID,
		OwnerID:        m.OwnerID,
		AvailableMoney: m.AvailableMoney,
		PendingMoney:   m.PendingMoney,
	}
	return wallet
}

func WalletModelListToWalletList(modelList []model.Wallet) []types.Wallet {
	walletList := make([]types.Wallet, len(modelList))
	for i, m := range modelList {
		wallet := WalletModelToWallet(m)
		walletList[i] = wallet
	}
	return walletList
}

// private

func (s *WalletService) takeWalletByModel(findModel *model.Wallet) (*gorm.DB, error) {
	tx := s.db.Where(*findModel).Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *WalletService) takeWalletBySeqAndModel(seq *gorm.DB, findModel *model.Wallet) (*gorm.DB, error) {
	tx := seq.Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

// private

func (s *WalletService) updateOneWalletByModel(findModel *model.Wallet, updateModel *model.Wallet) (*gorm.DB, error) {
	tx, err := s.takeOneWalletByModel(findModel)
	if err != nil {
		return nil, err
	}
	tx, err = s.updateOneWalletBySeqAndModel(tx, updateModel)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *WalletService) takeOneWalletByModel(findModel *model.Wallet) (*gorm.DB, error) {
	tx := s.db.Where(*findModel).Take(findModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (s *WalletService) updateOneWalletBySeqAndModel(seq *gorm.DB, updateModel *model.Wallet) (*gorm.DB, error) {
	tx := seq.Updates(updateModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}
