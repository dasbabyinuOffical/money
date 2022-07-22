package whale

import (
	"encoding/json"
	"fmt"
)

type Transaction struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Result  []Result `json:"result"`
}
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
	MethodID          string `json:"methodId"`
	FunctionName      string `json:"functionName"`
}

const (
	UrlPrefix = "https://api.etherscan.io/api"
	ApiKey    = "5Q5WXCCH8ZTSRKRJTY7K8XHZH5E5DT5T3S"
)

func GetTransactionByAddress(address string, startBlock uint64, endBlock uint64) (tx *Transaction, err error) {
	tx = new(Transaction)
	url := fmt.Sprintf("%s?module=account&action=tokentx&sort=asc&apiKey=%s&startblock=%d&endblock=%d&address=%s",
		UrlPrefix, ApiKey, startBlock, endBlock, address)
	respData, err := HttpGet(url)
	err = json.Unmarshal(respData, tx)
	if err != nil {
		return
	}
	return
}
