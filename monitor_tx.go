package main

import (
	"fmt"
	"money/whale"
	"time"
)

// TraceBscTx 1秒钟解析一次交易
func TraceBscTx() {
	for {
		time.Sleep(time.Second)
		txList, err := whale.FetchBscTransaction()
		if err != nil {
			fmt.Println("err is:", err)
		}
		for _, tx := range txList {
			security, score, err := whale.VerifyContract(tx)
			if err != nil {
				fmt.Println("VerifyContract err is:", err)
				continue
			}

			// 不安全的合约不用存
			if !security {
				continue
			}

			// 存储合约状态安全信息
			err = whale.SaveContractVerifyScore(score)
			if err != nil {
				fmt.Println("SaveContractVerifyScore,err is:", err)
				continue
			}

			// 存储所有交易
			err = whale.SaveBscTransaction(tx)
			if err != nil {
				fmt.Println("SaveBscTransaction,err is:", err)
				continue
			}

			// 存储热门币钟交易
			err = whale.SaveBsHotCoin(tx)
			if err != nil {
				fmt.Println("SaveBsHotCoin,err is:", err)
				continue
			}
		}
	}
}
