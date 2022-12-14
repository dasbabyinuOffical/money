package whale

import (
	"errors"
	"strconv"
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

func AnalysisContract(contract string, verify *ContractVerifyResult) (security bool, Score *ContractVerifyScore, err error) {
	// ???????????????
	if verify.Message != "OK" || len(verify.Result) == 0 {
		err = errors.New("???????????????")
		return
	}

	// ??????????????????
	var result *ContractVerify
	for _, value := range verify.Result {
		if value == nil {
			err = errors.New("???????????????nil")
			return
		}
		result = value
		break
	}

	// ??????????????????????????????
	if result.CannotSellAll == "1" {
		err = errors.New("?????????????????????")
		return
	}
	if result.CannotBuy == "1" {
		err = errors.New("???????????????")
		return
	}
	if result.CannotSellAll == "1" {
		err = errors.New("???????????????")
		return
	}
	if result.IsHoneypot == "1" {
		err = errors.New("????????????")
		return
	}
	if result.IsBlacklisted == "1" {
		err = errors.New("????????????")
		return
	}
	if result.IsMintable == "1" {
		err = errors.New("?????????")
		return
	}
	if result.IsOpenSource == "0" {
		err = errors.New("???????????????")
		return
	}
	if result.IsWhitelisted == "1" {
		err = errors.New("???????????????")
		return
	}
	if result.IsProxy == "1" {
		err = errors.New("???????????????")
		return
	}
	if result.TransferPausable == "1" {
		err = errors.New("???????????????")
		return
	}
	if result.TradingCooldown == "1" {
		err = errors.New("???????????????")
		return
	}
	if result.SlippageModifiable == "1" {
		err = errors.New("???????????????")
		return
	}

	// ????????????
	buyTax, err := strconv.ParseFloat(result.BuyTax, 10)
	if err != nil {
		return
	}
	// ??????????????????
	if buyTax >= 10 {
		err = errors.New("??????????????????")
		return
	}

	// ????????????
	sellTax, err := strconv.ParseFloat(result.SellTax, 10)
	if err != nil {
		return
	}
	// ??????????????????
	if sellTax >= 10 {
		err = errors.New("??????????????????")
		return
	}

	createrPercent, err := strconv.ParseFloat(result.CreatorPercent, 10)
	if err != nil {
		return
	}
	// ???????????????????????????
	if createrPercent >= 0.3 {
		return
	}

	// ?????????10????????????90%?????????????????????
	var top10Percent float64
	for _, holder := range result.Holders {
		percent, err := strconv.ParseFloat(holder.Percent, 10)
		if err != nil {
			return false, nil, err
		}
		top10Percent += percent
	}
	if top10Percent >= 90 {
		err = errors.New("???10????????????90%")
		return
	}

	var (
		lpLockPercent  float64
		lpTotalPercent float64
	)
	for _, lp := range result.LpHolders {
		lpPercent, err := strconv.ParseFloat(lp.Percent, 10)
		if err != nil {
			return false, nil, err
		}
		if lp.IsLocked == 1 && lp.IsContract == 1 {
			lpLockPercent += lpPercent
		}
		lpTotalPercent += lpPercent
	}
	if lpLockPercent <= 0 {
		err = errors.New("????????????")
		return
	}

	ownerPercent, err := strconv.ParseFloat(result.OwnerPercent, 10)
	if err != nil {
		return
	}

	//??????????????????
	totalSupply, err := strconv.ParseFloat(result.TotalSupply, 10)
	if err != nil {
		return
	}
	holderCount, err := strconv.ParseInt(result.HolderCount, 10, 64)
	if err != nil {
		return
	}

	// ???????????????????????????
	security = true
	Score = &ContractVerifyScore{
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
	if Score.BuyTax == 0 {
		Score.Score += 10
	}
	if Score.SellTax == 0 {
		Score.Score += 10
	}
	if Score.LPLockPercent >= 0.1 {
		Score.Score += 10
	}
	if Score.LpOfSupplyPercent >= 0.1 {
		Score.Score += 10
	}
	if Score.HolderCount <= 1000 {
		Score.Score += 10
	}
	// ???????????????
	if result.CreatorAddress != result.OwnerAddress {
		Score.Score += 10
		return
	}
	// ???????????????????????????50%
	if ownerPercent <= 0.5 || createrPercent <= 0.5 {
		Score.Score += 10
		return
	}
	// ???????????????0
	if buyTax == 0 || sellTax == 0 {
		Score.Score += 10
		return
	}
	// ???????????????10
	if buyTax >= 10 || sellTax >= 10 {
		Score.Score -= 10
		return
	}

	return
}
