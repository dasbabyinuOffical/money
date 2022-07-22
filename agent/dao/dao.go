package dao

import (
	"gorm.io/gorm"
	"money/agent/model"
	"time"
)

type Symbol string

var StableCoin = map[string]bool{
	"WBNB":    true,
	"BUSD":    true,
	"USDC":    true,
	"USDT":    true,
	"BNB":     true,
	"BSC-USD": true,
	"FIST":    true,
	"ETH":     true,
	"Cake":    true,
}

const (
	ContractVerifyUrlPrefix = "https://api.gopluslabs.io/api/v1/token_security/56?contract_addresses="
)

func SaveBscTransaction(tx *model.BSCTransaction) (err error) {
	today := time.Now().Format("2006-01-02")
	tx.Day = today
	tx.UpdatedAt = time.Now()

	// 如果买入卖出都是稳定币,则跳过不需要统计
	if StableCoin[tx.MakerSymbol] && StableCoin[tx.TakerSymbol] {
		return
	}

	oldTx := new(model.BSCTransaction)
	trx := DB().Begin()
	trx.Model(tx).
		Where("maker_contract = ? and taker_contract = ? and day = ?", tx.MakerContract, tx.TakerContract, today).
		First(oldTx)
	// 新增
	if oldTx.ID == 0 {
		tx.CreatedAt = time.Now()
		err = DB().Create(tx).Error
		trx.Commit()
		return
	}

	// 修改
	tx.TxCount = oldTx.TxCount + 1
	err = DB().Model(tx).Where("id = ?", oldTx.ID).Updates(tx).Error
	trx.Commit()
	return
}

func SaveBsHotCoin(tx *model.BSCTransaction) (err error) {
	// 如果买入卖出都是稳定币,则跳过不需要统计
	if StableCoin[tx.MakerSymbol] && StableCoin[tx.TakerSymbol] {
		return
	}

	hotCoin := new(model.HotCoin)
	today := time.Now()
	year, month, day, hour := today.Year(), today.Month(), today.Day(), today.Hour()
	hotCoin.Day = time.Date(year, month, day, hour, 0, 0, 0, time.Local)
	// 买入,maker是稳定币，取taker数据(卖出稳定币,购买新币)
	if _, ok := StableCoin[tx.MakerSymbol]; ok {
		oldTx := new(model.HotCoin)
		DB().Where("contract = ? and day = ?", tx.TakerContract, hotCoin.Day).First(oldTx)
		hotCoin.Symbol = tx.TakerSymbol
		hotCoin.Contract = tx.TakerContract
		hotCoin.TxCount = oldTx.TxCount + 1

		// 插入
		if oldTx.ID == 0 {
			err = DB().Create(hotCoin).Error
			return
		}
		// 修改
		err = DB().Model(hotCoin).Where("id = ?", oldTx.ID).Updates(hotCoin).Error
		return
	}

	// 卖出,taker是稳定币，取maker数据(卖出新币,购买稳定币)
	oldTx := new(model.HotCoin)
	DB().Where("contract = ? and day = ?", tx.MakerContract, hotCoin.Day).First(oldTx)
	hotCoin.Symbol = tx.MakerSymbol
	hotCoin.Contract = tx.MakerContract
	hotCoin.TxCount = oldTx.TxCount + 1

	// 插入
	if oldTx.ID == 0 {
		err = DB().Create(hotCoin).Error
		return
	}
	// 修改
	err = DB().Model(hotCoin).Where("id = ?", oldTx.ID).Updates(hotCoin).Error
	return
}

func GetBscTransaction(makerContract string, takerContract string, day string) (bscTx *model.BSCTransaction, err error) {
	bscTx = new(model.BSCTransaction)
	err = DB().Where("maker_contract = ? and taker_contract = ? and day = ?",
		makerContract, takerContract, day).
		First(bscTx).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

func GetBscLatestTransactionFromDB() (txList []*model.BSCTransaction, err error) {
	err = DB().Order("updated_at desc").
		Limit(100).Find(&txList).Error
	return
}

func GetBscHotTransactionFromDB() (hotCoinList []*model.HotCoin, err error) {
	today := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	err = DB().Model(new(model.HotCoin)).Joins("inner join contract_verify_score on hot_coin.contract = contract_verify_score.contract").
		Where("hot_coin.day >= ?", today).
		Order("contract_verify_score.score desc,hot_coin.tx_count desc").
		Limit(100).Find(&hotCoinList).Error
	return
}

func GetContractVerifyScore(contract string) (score *model.ContractVerifyScore, err error) {
	score = new(model.ContractVerifyScore)
	err = DB().Where("contract = ?", contract).First(score).Error
	return
}

func SaveContractVerifyScore(score *model.ContractVerifyScore) (err error) {
	err = DB().Where("contract = ?", score.Contract).FirstOrCreate(score).Error
	return
}

func GetAddressBlock(address string) (addressBlock *model.AddressBlock, err error) {
	addressBlock = new(model.AddressBlock)
	err = DB().Model(new(model.AddressBlock)).
		Where("address = ?", address).
		First(addressBlock).Error
	return
}

func SaveAddressBlock(addressBlock *model.AddressBlock) (err error) {
	err = DB().Where("id = ?", addressBlock.ID).Updates(addressBlock).Error
	return
}
