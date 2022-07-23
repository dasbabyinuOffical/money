package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"money/agent/contracts"
	"money/agent/contracts/pancakev2"
	"money/agent/dao"
	"money/agent/model"
	"money/agent/tools"
	"strconv"
	"time"
)

type Agent struct {
	PancakeV2Caller *pancakev2.Pancakev2Caller
	PancakeV2ABI    string
	BscClient       *ethclient.Client
}

type ContractPrice struct {
	UpdatedAt int64 `json:"updated_at"`
	Data      struct {
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		Price    string `json:"price"`
		PriceBNB string `json:"price_BNB"`
	} `json:"data"`
}

func NewAgent() (agent *Agent, err error) {
	agent = new(Agent)

	// 获取bscClient
	agent.BscClient, err = ethclient.Dial(contracts.BscRpcUrl)
	if err != nil {
		return
	}

	// 获取pancakeV2Caller
	agent.PancakeV2Caller, err = contracts.NewPancakeV2PoolCaller(contracts.BscRpcUrl)
	if err != nil {
		return
	}

	// 获取pancakeV2ABI
	abiUrl := fmt.Sprintf("%s?module=contract&action=getabi&address=%s&&apikey=%s",
		contracts.BscScanUrl, contracts.PancakeV2Pool, contracts.BscScanAPIKEY)
	data, err := tools.HttpGet(abiUrl)
	if err != nil {
		return
	}
	agent.PancakeV2ABI = string(data)

	return
}

//ScanLog 扫描区块日志
func (agent *Agent) ScanLog() (err error) {
	// 获取薄饼路由最新扫描到的区块编号,加入新合约需要手动插入区块编号
	addressBlock, err := dao.GetAddressBlock(contracts.PancakeV2Pool)
	if err != nil {
		return
	}

	// 获取链上最新的区块编号
	blockNumber, err := agent.BscClient.BlockNumber(context.Background())
	if err != nil || blockNumber == 0 {
		return
	}
	blockNumberOnChain := big.NewInt(int64(blockNumber))
	blockNumberInDB := new(big.Int)
	blockNumberInDB.SetInt64(addressBlock.Block)
	if blockNumberOnChain.Cmp(blockNumberInDB) < 1 {
		log.Println("链上区块小于当前区块，无需扫描事件")
		return
	}

	txUrl := fmt.Sprintf("%s?module=account&action=txlist&sort=asc&apiKey=%s&startblock=%d&endblock=%d&address=%s",
		contracts.BscScanUrl, contracts.BscScanAPIKEY, blockNumberInDB.Int64(), blockNumberOnChain.Int64(), contracts.PancakeV2Pool)
	// 获取最新区块日志
	txRecord, err := dao.FetchBlockResults(txUrl)
	if err != nil || txRecord.Status != "1" || txRecord.Message != "OK" {
		return
	}
	for _, result := range txRecord.Result {
		// 交易错误或者挂起交易，进行跳过
		if result.IsError != "0" || result.TxreceiptStatus != "1" {
			fmt.Println(result.Hash, ":交易记录错误或者receipt状态错误")
			continue
		}
		record, err := agent.AnalysisTxRecord(result)
		if err != nil {
			fmt.Println(record, err)
			continue
		}

		// 如果买入卖出都是稳定币,则跳过不需要统计
		if dao.IsStableCoin(record.MakerContract) && dao.IsStableCoin(record.TakerContract) {
			fmt.Println("稳定币交易跳过")
			continue
		}

		// 分析交易，验证合约安全性
		scores, err := dao.VerifyContract(record)
		if err != nil || len(scores) == 0 {
			fmt.Println("合约不安全:", record.TxId)
			continue
		}
		security := true
		for _, score := range scores {
			security = security && score.Security
		}
		if security == false {
			fmt.Println("合约不安全,跳过:", record.TxId)
			continue
		}

		// 获取合约价格和symbol信息
		var makerPriceInfo ContractPrice
		data, err := tools.HttpGet(contracts.PancakeV2Route + record.MakerContract)
		if err != nil {
			fmt.Println("获取合约价格失败," + record.MakerContract)
			continue
		}
		err = json.Unmarshal(data, &makerPriceInfo)
		if err != nil {
			fmt.Println("合约价格反序列化失败," + record.MakerContract)
			continue
		}
		record.MakerSymbol = makerPriceInfo.Data.Symbol
		record.MakerPrice, _ = strconv.ParseFloat(makerPriceInfo.Data.Price, 10)

		var takerPriceInfo ContractPrice
		data, err = tools.HttpGet(contracts.PancakeV2Route + record.TakerContract)
		if err != nil {
			fmt.Println("获取合约价格失败," + record.TakerContract)
			continue
		}
		err = json.Unmarshal(data, &takerPriceInfo)
		if err != nil {
			fmt.Println("合约价格反序列化失败," + record.TakerContract)
			continue
		}
		record.TakerSymbol = takerPriceInfo.Data.Symbol
		record.TakerPrice, _ = strconv.ParseFloat(takerPriceInfo.Data.Price, 10)

		// 保存合约
		for _, score := range scores {
			err = dao.SaveContractVerifyScore(score)
			if err != nil {
				fmt.Println(" 保存验证合约失败：", score.Contract)
				continue
			}
		}

		// 保存交易记录
		err = dao.SaveBscTransaction(record)
		if err != nil {
			fmt.Println("交易记录保存失败:", record.TxId)
			continue
		}

		// 保存热门币合约
		hotCoins, err := dao.TransferToBscHotCoin(record)
		if err != nil {
			fmt.Println("热门币获取失败:", record.TxId)
			continue
		}
		for _, hotCoin := range hotCoins {
			dao.SaveHotCoin(hotCoin)
		}
	}
	// 更新合约扫描区块编号
	addressBlock.Block = blockNumberOnChain.Int64()
	err = dao.SaveAddressBlock(addressBlock)
	return
}

// 处理区块日志
func (agent *Agent) AnalysisTxRecord(txRecord *dao.Result) (bscTx *model.BSCTransaction, err error) {
	bscTx = new(model.BSCTransaction)

	// 解析交易输入和输出
	decoder, err := NewBscInputDecoder(txRecord.FunctionName)
	if err != nil {
		return
	}

	inputTx, err := decoder.DecodeInput(txRecord.Input)
	if err != nil {
		return
	}

	timeStamp, err := strconv.ParseInt(txRecord.TimeStamp, 10, 64)
	if err != nil {
		return
	}
	txTime := time.Unix(timeStamp, 0)
	day := txTime.Format("2006-01-02")

	// 从数据库获取交易对
	bscTx, err = dao.GetBscTransaction(inputTx.MakerContract, inputTx.TakerContract, day)
	if err != nil {
		return
	}

	bscTx.TxId = txRecord.Hash
	bscTx.TxTime = txTime
	bscTx.Day = day
	bscTx.TxCount += 1
	bscTx.MakerContract = inputTx.MakerContract
	bscTx.TakerContract = inputTx.TakerContract
	return
}

// 定时更新合约相关信息
func (agent *Agent) SyncContractInfo() (er error) {
	// 获取所有未更新合约时间的合约
	scores, err := dao.GetContractUnCreatedList()
	if err != nil {
		return
	}
	for _, score := range scores {
		// 现在api频率
		time.Sleep(time.Second)
		// 获取合约创建时间
		txUrl := fmt.Sprintf("%s?module=account&action=txlist&sort=asc&apiKey=%s&startblock=%d&page=1&offset=1&address=%s",
			contracts.BscScanUrl, contracts.BscScanAPIKEY2, 0, score.Contract)
		// 获取最新区块日志
		txRecord, err := dao.FetchBlockResults(txUrl)
		if err != nil || txRecord.Status != "1" || txRecord.Message != "OK" || len(txRecord.Result) == 0 {
			return
		}
		ts, err := strconv.ParseInt(txRecord.Result[0].TimeStamp, 10, 64)
		if err != nil {
			return
		}
		createdDay := time.Unix(ts, 0)
		score.CreatedDay = &createdDay
		dao.SaveContractVerifyScore(score)
	}
	return
}
