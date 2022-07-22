package dao

import (
	"gorm.io/gorm"
	"money/agent/model"
	"strings"
	"time"
)

type Symbol string

var StableCoin = map[string]string{
	"0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c": "WBNB",
	"0xe9e7cea3dedca5984780bafc599bd69add087d56": "BUSD",
	"0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d": "USDC",
	"0x55d398326f99059ff775485246999027b3197955": "BSC-USD",
	"0x2170ed0880ac9a755fd29b2688956bd959f933f8": "ETH",
	"0x0e09fabb73bd3ade0a17ecc321fd13a19e81ce82": "Cake",
}

const (
	ContractVerifyUrlPrefix = "https://api.gopluslabs.io/api/v1/token_security/56?contract_addresses="
)

func SaveBscTransaction(tx *model.BSCTransaction) (err error) {
	today := time.Now().Format("2006-01-02")
	tx.Day = today
	tx.UpdatedAt = time.Now()

	oldTx := new(model.BSCTransaction)
	DB().Model(tx).
		Where("maker_contract = ? and taker_contract = ? and day = ?", tx.MakerContract, tx.TakerContract, today).
		First(oldTx)
	// 新增
	if oldTx.ID == 0 {
		tx.CreatedAt = time.Now()
		err = DB().Create(tx).Error
		return
	}

	// 修改
	tx.TxCount = oldTx.TxCount + 1
	err = DB().Model(tx).Where("id = ?", oldTx.ID).Updates(tx).Error
	return
}

func TransferToBscHotCoin(tx *model.BSCTransaction) (hotCoins []*model.HotCoin, err error) {
	today := time.Now()
	year, month, day, hour := today.Year(), today.Month(), today.Day(), today.Hour()

	if !IsStableCoin(tx.MakerContract) {
		hotCoin := new(model.HotCoin)
		hotCoin.Day = time.Date(year, month, day, hour, 0, 0, 0, time.Local)
		oldTx := new(model.HotCoin)
		DB().Where("contract = ? and day = ?", tx.MakerContract, hotCoin.Day).First(oldTx)
		hotCoin.Symbol = tx.MakerSymbol
		hotCoin.Contract = tx.MakerContract
		hotCoin.TxCount = oldTx.TxCount + 1
		hotCoin.Price = tx.MakerPrice
		hotCoins = append(hotCoins, hotCoin)
	}

	if !IsStableCoin(tx.TakerContract) {
		hotCoin := new(model.HotCoin)
		hotCoin.Day = time.Date(year, month, day, hour, 0, 0, 0, time.Local)
		oldTx := new(model.HotCoin)
		DB().Where("contract = ? and day = ?", tx.TakerContract, hotCoin.Day).First(oldTx)
		hotCoin.Symbol = tx.TakerSymbol
		hotCoin.Contract = tx.TakerContract
		hotCoin.TxCount = oldTx.TxCount + 1
		hotCoin.Price = tx.TakerPrice
		hotCoins = append(hotCoins, hotCoin)
	}
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
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

func SaveContractVerifyScore(score *model.ContractVerifyScore) (err error) {
	if score.ID == 0 {
		err = DB().Create(score).Error
		return
	}
	err = DB().Where("id = ?", score.ID).Updates(score).Error
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

func SaveHotCoin(hotCoin *model.HotCoin) (err error) {
	if hotCoin.ID == 0 {
		err = DB().Create(hotCoin).Error
		return
	}
	err = DB().Where("id = ?", hotCoin.ID).Updates(hotCoin).Error
	return
}

func IsStableCoin(contract string) bool {
	_, ok := StableCoin[strings.ToLower(contract)]
	return ok
}