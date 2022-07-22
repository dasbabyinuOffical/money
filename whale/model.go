package whale

import (
	"gorm.io/gorm"
	"time"
)

type AddressTransaction struct {
	ID        uint64         `gorm:"column:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`
	Address   string         `gorm:"column:address" json:"address"`
	TxId      string         `gorm:"column:tx_id" json:"txId"`
	Block     uint64         `gorm:"column:block" json:"block"`
}

func (t *AddressTransaction) TableName() string {
	return "address_transaction"
}

type BSCTransaction struct {
	ID            uint64         `gorm:"column:id" json:"id"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`
	TxId          string         `gorm:"column:tx_id" json:"txId"`
	TxTime        time.Time      `gorm:"column:tx_time" json:"txTime"`
	Time          string         `gorm:"-"`
	MakerAmount   float64        `gorm:"column:maker_amount" json:"makerAmount"`
	MakerSymbol   string         `gorm:"column:maker_symbol" json:"makerSymbol"`
	MakerContract string         `gorm:"column:maker_contract" json:"makerContract"`
	TakerAmount   float64        `gorm:"column:taker_amount" json:"takerAmount"`
	TakerSymbol   string         `gorm:"column:taker_symbol" json:"takerSymbol"`
	TakerContract string         `gorm:"column:taker_contract" json:"takerContract"`
	Price         float64        `gorm:"column:price" json:"price"`
	PriceSymbol   string         `gorm:"column:price_symbol" json:"priceSymbol"`
	Status        uint8          `gorm:"column:status" json:"status"`
	Day           string         `gorm:"column:day" json:"day"`
	TxCount       uint64         `gorm:"column:tx_count" json:"txCount"`
	MakerVolume   float64        `gorm:"column:maker_volume" json:"makerVolume"`
	TakerVolume   float64        `gorm:"column:taker_volume" json:"takerVolume"`
}

func (t *BSCTransaction) TableName() string {
	return "bsc_transaction"
}

// HotCoin 今日热门币,前10
type HotCoin struct {
	ID          uint64         `gorm:"column:id" json:"id"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`
	Contract    string         `gorm:"column:contract" json:"contract"`
	Symbol      string         `gorm:"column:symbol" json:"symbol"`
	Price       float64        `gorm:"column:price" json:"price"`
	PriceSymbol string         `gorm:"column:price_symbol" json:"priceSymbol"`
	MakerVolume float64        `gorm:"column:maker_volume" json:"makerVolume"`
	TakerVolume float64        `gorm:"column:taker_volume" json:"takerVolume"`
	TotalVolume float64        `gorm:"column:total_volume" json:"totalVolume"`
	TxCount     uint64         `gorm:"column:tx_count" json:"txCount"`
	Day         time.Time      `gorm:"column:day" json:"day"`
}

func (t *HotCoin) TableName() string {
	return "hot_coin"
}

type ContractVerifyScore struct {
	ID          uint64         `gorm:"column:id" json:"id"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt"`
	TokenName   string         `gorm:"column:token_name" json:"token_name"`
	TokenSymbol string         `gorm:"column:token_symbol" json:"token_symbol"`
	Contract    string         `gorm:"column:contract" json:"contract"`
	TotalSupply float64        `gorm:"column:total_supply" json:"total_supply"`

	// 创建人信息
	CreatorAddress string  `gorm:"column:creator_address" json:"creator_address"`
	CreatorPercent float64 `gorm:"column:creator_percent" json:"creator_percent"`
	CreatedDay     string  `gorm:"column:created_day" json:"created_day"`

	// 合约所有者信息
	OwnerAddress       string  `gorm:"column:owner_address" json:"owner_address"`
	OwnerPercent       float64 `gorm:"column:owner_percent" json:"owner_percent"`
	HolderCount        int64   `gorm:"column:holder_count" json:"holder_count"`
	Top10HolderPercent float64 `gorm:"column:top_10_holder_percent" json:"top_10_holder_percent"`

	//LP信息
	LPLockPercent     float64 `gorm:"column:lp_lock_percent" json:"lp_lock_percent"`
	LpOfSupplyPercent float64 `gorm:"column:lp_of_supply_percent" json:"lp_of_supply_percent"`

	// 安全信息
	BuyTax  float64 `gorm:"column:buy_tax" json:"buy_tax"`
	SellTax float64 `gorm:"column:sell_tax" json:"sell_tax"`

	Circulation        float64 `gorm:"column:circulation" json:"circulation"`
	CirculationPercent float64 `gorm:"column:circulation_percent" json:"circulation_percent"`
	MarketCap          uint64  `gorm:"column:market_cap" json:"market_cap"`
	CurrentPrice       float64 `gorm:"column:current_price" json:"current_price"`
	PriceUp            uint64  `gorm:"column:price_up" json:"price_up"`     //距离开盘价上涨了多少倍
	PriceDown          uint64  `gorm:"column:price_down" json:"price_down"` //距离最高价跌了多少倍
	Website            string  `gorm:"column:website" json:"website"`
	Telegram           string  `gorm:"column:telegram" json:"telegram"`
	Twitter            string  `gorm:"column:twitter" json:"twitter"`
	Discord            string  `gorm:"column:discord" json:"discord"`
	Youtube            string  `gorm:"column:youtube" json:"youtube"`
	TxCount            uint64  `gorm:"column:tx_count" json:"tx_count"`
	// 综合得分以及是否安全
	Score    int64 `gorm:"column:score" json:"score"`
	Security bool  `gorm:"column:security" json:"security"`
}

func (t *ContractVerifyScore) TableName() string {
	return "contract_verify_score"
}
