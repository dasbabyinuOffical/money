package dao

import (
	"encoding/json"
	"errors"
	"money/agent/model"
	"money/agent/tools"
	"strconv"
	"strings"
)

type ContractVerifyResult struct {
	Code    int                        `json:"code"`
	Message string                     `json:"message"`
	Result  map[string]*ContractVerify `json:"result"`
}
type Dex struct {
	Name      string `json:"name"`
	Liquidity string `json:"liquidity"`
	Pair      string `json:"pair"`
}
type Holders struct {
	Address    string `json:"address"`
	Tag        string `json:"tag"`
	IsContract int    `json:"is_contract"`
	Balance    string `json:"balance"`
	Percent    string `json:"percent"`
	IsLocked   int    `json:"is_locked"`
}
type LpHolders struct {
	Address    string `json:"address"`
	Tag        string `json:"tag"`
	IsContract int    `json:"is_contract"`
	Balance    string `json:"balance"`
	Percent    string `json:"percent"`
	IsLocked   int    `json:"is_locked"`
}
type ContractVerify struct {
	BuyTax               string      `json:"buy_tax"`
	CanTakeBackOwnership string      `json:"can_take_back_ownership"`
	CannotBuy            string      `json:"cannot_buy"`
	CannotSellAll        string      `json:"cannot_sell_all"`
	CreatorAddress       string      `json:"creator_address"`
	CreatorBalance       string      `json:"creator_balance"`
	CreatorPercent       string      `json:"creator_percent"`
	Dex                  []Dex       `json:"dex"`
	HiddenOwner          string      `json:"hidden_owner"`
	HolderCount          string      `json:"holder_count"`
	Holders              []Holders   `json:"holders"`
	IsAntiWhale          string      `json:"is_anti_whale"`
	IsBlacklisted        string      `json:"is_blacklisted"`
	IsHoneypot           string      `json:"is_honeypot"`
	IsInDex              string      `json:"is_in_dex"`
	IsMintable           string      `json:"is_mintable"`
	IsOpenSource         string      `json:"is_open_source"`
	IsProxy              string      `json:"is_proxy"`
	IsWhitelisted        string      `json:"is_whitelisted"`
	LpHolderCount        string      `json:"lp_holder_count"`
	LpHolders            []LpHolders `json:"lp_holders"`
	LpTotalSupply        string      `json:"lp_total_supply"`
	OwnerAddress         string      `json:"owner_address"`
	OwnerBalance         string      `json:"owner_balance"`
	OwnerChangeBalance   string      `json:"owner_change_balance"`
	OwnerPercent         string      `json:"owner_percent"`
	SellTax              string      `json:"sell_tax"`
	SlippageModifiable   string      `json:"slippage_modifiable"`
	TotalSupply          string      `json:"total_supply"`
	TradingCooldown      string      `json:"trading_cooldown"`
	TransferPausable     string      `json:"transfer_pausable"`
	TokenName            string      `json:"token_name"`
	TokenSymbol          string      `json:"token_symbol"`
}

func verifyContract(contract string) (score *model.ContractVerifyScore, err error) {
	score = new(model.ContractVerifyScore)
	score, err = GetContractVerifyScore(contract)
	// 验证合约已经存在或取不到则不用再验证
	if err != nil || score.ID > 0 {
		return
	}

	url := ContractVerifyUrlPrefix + contract
	data, err := tools.HttpGet(url)
	if err != nil {
		return
	}

	var result = new(ContractVerifyResult)
	err = json.Unmarshal(data, result)
	if err != nil {
		return
	}
	score, err = AnalysisContract(contract, result)
	if err != nil {
		return
	}
	return
}

func VerifyContract(tx *model.BSCTransaction) (scores []*model.ContractVerifyScore, err error) {
	if !IsStableCoin(tx.MakerContract) {
		if score, err := verifyContract(tx.MakerContract); err == nil {
			scores = append(scores, score)
		}
	}
	if !IsStableCoin(tx.TakerContract) {
		if score, err := verifyContract(tx.TakerContract); err == nil {
			scores = append(scores, score)
		}
	}
	return
}

func contractEqual(contract1, contract2 string) bool {
	return strings.ToLower(contract1) == strings.ToLower(contract2)
}

func AnalysisContract(contract string, verify *ContractVerifyResult) (score *model.ContractVerifyScore, err error) {
	// 数据无返回
	if verify.Message != "OK" || len(verify.Result) == 0 {
		err = errors.New("未返回数据")
		return
	}

	// 获取验证数据
	var result *ContractVerify
	for key, value := range verify.Result {
		if contractEqual(key, contract) && value != nil {
			result = value
			break
		}
	}
	if result == nil {
		err = errors.New("返回合约验证数据错误")
		return
	}

	// 判断不安全项直接返回
	if result.CannotSellAll == "1" {
		err = errors.New("不允许全部卖出")
		return
	}
	if result.CannotBuy == "1" {
		err = errors.New("不允许买入")
		return
	}
	if result.CannotSellAll == "1" {
		err = errors.New("不允许卖出")
		return
	}
	if result.IsHoneypot == "1" {
		err = errors.New("包含蜜罐")
		return
	}
	if result.IsBlacklisted == "1" {
		err = errors.New("有黑名单")
		return
	}

	if result.IsWhitelisted == "1" {
		err = errors.New("包含白名单")
		return
	}

	if result.TransferPausable == "1" {
		err = errors.New("交易可暂停")
		return
	}
	if result.TradingCooldown == "1" {
		err = errors.New("交易可关闭")
		return
	}

	// 买入税率
	buyTax, err := strconv.ParseFloat(result.BuyTax, 10)
	if err != nil {
		return
	}
	// 买入税率太高
	if buyTax >= 0.2 {
		err = errors.New("买入税率太高")
		return
	}

	// 卖出税率
	sellTax, err := strconv.ParseFloat(result.SellTax, 10)
	if err != nil {
		return
	}
	// 卖出税率太高
	if sellTax >= 0.2 {
		err = errors.New("卖出税率太高")
		return
	}

	createrPercent, err := strconv.ParseFloat(result.CreatorPercent, 10)
	if err != nil {
		return
	}
	// 创建者拥有太多的币
	if createrPercent >= 0.3 {
		err = errors.New("创建人拥有太多的币")
		return
	}

	// 持币前10是否超过99%的币，持币太多
	var top10Percent float64
	for _, holder := range result.Holders {
		percent, err := strconv.ParseFloat(holder.Percent, 10)
		if err != nil {
			return nil, err
		}
		top10Percent += percent
	}
	if top10Percent >= 0.99 {
		err = errors.New("前10持币超过100%")
		return
	}

	var (
		lpLockPercent  float64
		lpTotalPercent float64
	)
	for _, lp := range result.LpHolders {
		lpPercent, err := strconv.ParseFloat(lp.Percent, 10)
		if err != nil {
			return nil, err
		}
		if lp.IsLocked == 1 || lp.IsContract == 1 {
			lpLockPercent += lpPercent
		}
		lpTotalPercent += lpPercent
	}
	if lpLockPercent <= 0 {
		err = errors.New("池子未锁")
		return
	}

	ownerPercent, err := strconv.ParseFloat(result.OwnerPercent, 10)
	if err != nil {
		return
	}

	//计算结果数据
	totalSupply, err := strconv.ParseFloat(result.TotalSupply, 10)
	if err != nil {
		return
	}
	holderCount, err := strconv.ParseInt(result.HolderCount, 10, 64)
	if err != nil {
		return
	}

	// 到这里合约才算安全
	score = &model.ContractVerifyScore{
		TokenName:          result.TokenName,
		TokenSymbol:        result.TokenSymbol,
		Contract:           contract,
		TotalSupply:        totalSupply,
		CreatorAddress:     result.CreatorAddress,
		CreatorPercent:     createrPercent,
		OwnerAddress:       result.OwnerAddress,
		OwnerPercent:       ownerPercent,
		HolderCount:        holderCount,
		Top10HolderPercent: top10Percent,
		LPLockPercent:      lpLockPercent,
		LpOfSupplyPercent:  lpTotalPercent,
		BuyTax:             buyTax,
		SellTax:            sellTax,
		Score:              0,
		Security:           true,
	}
	if score.BuyTax == 0 {
		score.Score += 10
	}
	if score.SellTax == 0 {
		score.Score += 10
	}
	if score.LPLockPercent >= 0.1 {
		score.Score += 10
	}
	if score.LpOfSupplyPercent >= 0.1 {
		score.Score += 10
	}

	// 新盘或久经考验的老盘
	if holderCount <= 1000 || holderCount >= 10000 {
		score.Score += 20
	}
	// 放弃所有权
	if result.CreatorAddress != result.OwnerAddress {
		score.Score += 10
	}
	// 合约所有者比率小于50%
	if ownerPercent <= 0.5 || createrPercent <= 0.5 {
		score.Score += 10
	}
	// 买卖税均为0
	if buyTax == 0 || sellTax == 0 {
		score.Score += 10
	}
	// 买卖税大于10
	if buyTax >= 10 || sellTax >= 10 {
		score.Score -= 10
	}

	// 可增发
	if result.IsMintable == "1" {
		score.Score -= 10
	}
	// 合约未开源
	if result.IsOpenSource == "0" {
		score.Score -= 10
	}
	// 代理合约
	if result.IsProxy == "1" {
		score.Score -= 10
	}
	// 可以修改交易费率
	if result.SlippageModifiable == "1" {
		score.Score -= 10
	}
	return
}
