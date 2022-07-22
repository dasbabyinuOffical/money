package dao

import (
	"encoding/json"
	"money/agent/tools"
)

type Result struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxreceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
	FunctionName      string `json:"functionName"`
}

type TxRecord struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Result  []*Result `json:"result"`
}

func FetchBlockResults(url string) (txRecord *TxRecord, err error) {
	txRecord = new(TxRecord)
	respData, err := tools.HttpGet(url)
	err = json.Unmarshal(respData, txRecord)
	if err != nil {
		return
	}
	return
}
