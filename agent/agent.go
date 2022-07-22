package agent

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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

const (
	BscRpcUrl     = "https://bsc-dataseed1.binance.org"
	BscScanUrl    = "https://api.bscscan.com/api"
	BscScanAPIKEY = "FT1ZFNFUJV1GRDKZQB4FXDUYZD31CVAIWE"
	Swap          = "Swap (index_topic_1 address sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, index_topic_2 address to)"
)

type Agent struct {
	PancakeV2Caller *pancakev2.Pancakev2Caller
	PancakeV2ABI    string
	BscClient       *ethclient.Client
}

func NewAgent() (agent *Agent, err error) {
	agent = new(Agent)

	// 获取bscClient
	agent.BscClient, err = ethclient.Dial(BscRpcUrl)
	if err != nil {
		return
	}

	// 获取pancakeV2Caller
	agent.PancakeV2Caller, err = contracts.NewPancakeV2PoolCaller(BscRpcUrl)
	if err != nil {
		return
	}

	// 获取pancakeV2ABI
	abiUrl := fmt.Sprintf("%s?module=contract&action=getabi&address=%s&&apikey=%s",
		BscScanUrl, contracts.PancakeV2Pool, BscScanAPIKEY)
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
		BscScanUrl, BscScanAPIKEY, blockNumberInDB.Int64(), blockNumberOnChain.Int64(), contracts.PancakeV2Pool)
	// 获取最新区块日志
	txRecord, err := fetchBlockResults(txUrl)
	if err != nil || txRecord.Status != "1" || txRecord.Message != "OK" {
		return
	}
	for _, result := range txRecord.Result {
		record, err := agent.AnalysisTxRecord(result)
		if err != nil {
			fmt.Println(record, err)
			continue
		}

		// 分析交易，验证合约安全性
		score, err := dao.VerifyContract(record)
		if err != nil {
			fmt.Println("合约无法验证:", record.TxId)
			continue
		}
		if score.Security == false {
			fmt.Println("合约不安全，跳过:", record.TxId)
			continue
		}

		// 保存合约
		err = dao.SaveContractVerifyScore(score)
		if err != nil {
			fmt.Println(" 保存验证合约失败：", score.Contract)
			continue
		}

		// 保存交易记录
		err = dao.SaveBscTransaction(record)
		if err != nil {
			fmt.Println("交易记录保存失败:", record.TxId)
			continue
		}

		// 保存热门币合约
		err = dao.SaveBsHotCoin(record)
		if err != nil {
			fmt.Println("热门币保存失败:", record.TxId)
			continue
		}
	}
	// 更新合约扫描区块编号
	addressBlock.Block = blockNumberOnChain.Int64()
	dao.SaveAddressBlock(addressBlock)
	return
}

// 处理区块日志
func (agent *Agent) AnalysisTxRecord(txRecord *Result) (bscTxResult *model.BSCTransaction, err error) {
	// 交易错误或者挂起交易，进行跳过
	if txRecord.IsError != "0" || txRecord.TxreceiptStatus != "1" {
		return
	}

	// 获取交易日志
	receipt, err := agent.BscClient.TransactionReceipt(context.Background(), common.HexToHash(txRecord.Hash))
	if err != nil {
		return
	}
	var logData []byte
	for _, log := range receipt.Logs {
		method := log.Topics[0].Hex()
		if tools.SignABI(method) == Swap {
			logData = log.Data
		}
	}

	// 解析交易输入和输出
	decoder, err := NewBscInputDecoder(txRecord.FunctionName)
	if err != nil {
		return
	}

	inputTx, err := decoder.DecodeInput(txRecord.Input, logData)
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
	bscTx, err := dao.GetBscTransaction(inputTx.MakerContract, inputTx.TakerContract, day)
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
