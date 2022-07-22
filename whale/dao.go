package whale

import (
	"gorm.io/gorm"
	"time"
)

func FetchAddressLastBlock(address string) (blockNum uint64, err error) {
	addressTransaction := new(AddressTransaction)

	err = DB().Model(addressTransaction).
		Where("address = ?", address).
		Last(addressTransaction).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	blockNum = addressTransaction.Block
	return
}

func SaveBscTransaction(tx *BSCTransaction) (err error) {
	today := time.Now().Format("2006-01-02")
	tx.Day = today
	tx.UpdatedAt = time.Now()

	// 如果买入卖出都是稳定币,则跳过不需要统计
	if StableCoin[tx.MakerSymbol] && StableCoin[tx.TakerSymbol] {
		return
	}

	oldTx := new(BSCTransaction)
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
	if _, ok := StableCoin[tx.MakerSymbol]; ok {
		tx.MakerVolume = oldTx.MakerVolume + tx.MakerAmount
	}
	if _, ok := StableCoin[tx.TakerSymbol]; ok {
		tx.TakerVolume = oldTx.TakerVolume + tx.TakerAmount
	}
	tx.TxCount = oldTx.TxCount + 1
	err = DB().Model(tx).Where("id = ?", oldTx.ID).Updates(tx).Error
	trx.Commit()
	return
}

func SaveBsHotCoin(tx *BSCTransaction) (err error) {
	// 如果买入卖出都是稳定币,则跳过不需要统计
	if StableCoin[tx.MakerSymbol] && StableCoin[tx.TakerSymbol] {
		return
	}

	hotCoin := new(HotCoin)
	today := time.Now()
	year, month, day, hour := today.Year(), today.Month(), today.Day(), today.Hour()
	hotCoin.Day = time.Date(year, month, day, hour, 0, 0, 0, time.Local)
	// 买入,maker是稳定币，取taker数据(卖出稳定币,购买新币)
	if _, ok := StableCoin[tx.MakerSymbol]; ok {
		oldTx := new(HotCoin)
		DB().Where("contract = ? and day = ?", tx.TakerContract, hotCoin.Day).First(oldTx)
		hotCoin.Symbol = tx.TakerSymbol
		hotCoin.Contract = tx.TakerContract
		hotCoin.Price = tx.Price
		hotCoin.PriceSymbol = tx.PriceSymbol
		hotCoin.TakerVolume = oldTx.TakerVolume + tx.TakerAmount
		hotCoin.TotalVolume = oldTx.TotalVolume + tx.TakerAmount
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
	oldTx := new(HotCoin)
	DB().Where("contract = ? and day = ?", tx.MakerContract, hotCoin.Day).First(oldTx)
	hotCoin.Symbol = tx.MakerSymbol
	hotCoin.Contract = tx.MakerContract
	hotCoin.MakerVolume = oldTx.MakerVolume + tx.MakerAmount
	hotCoin.TotalVolume = oldTx.TotalVolume + tx.MakerAmount
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

func GetBscLatestTransactionFromDB() (txList []*BSCTransaction, err error) {
	err = DB().Order("updated_at desc").
		Limit(100).Find(&txList).Error
	return
}

func GetBscHotTransactionFromDB() (hotCoinList []*HotCoin, err error) {
	today := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	err = DB().Model(new(HotCoin)).Joins("inner join contract_verify_score on hot_coin.contract = contract_verify_score.contract").
		Where("hot_coin.day >= ?", today).
		Order("contract_verify_score.score desc,hot_coin.tx_count desc").
		Limit(100).Find(&hotCoinList).Error
	return
}

func IsContractHasVerified(contract string) (flag bool) {
	score := new(ContractVerifyScore)
	err = DB().Where("contract = ?", contract).First(score).Error
	if err == nil && score.ID > 0 {
		return true
	}
	return
}

func SaveContractVerifyScore(score *ContractVerifyScore) (err error) {
	err = DB().Where("contract = ?", score.Contract).FirstOrCreate(score).Error
	return
}
