package main

import (
	"log"
	"money/whale"
	"time"
)

var (
	addressList = []string{
		"0x11da0d7401319585ca421a3715888b362da7fe6b",
	}
	to = []string{
		"394329020@qq.com",
		"15228328910@139.com",
	}
)

func scanTx(address string, endBlock uint64) (results []whale.Result, err error) {
	startBlock, err := whale.FetchAddressLastBlock(address)
	if err != nil {
		return
	}

	tx, err := whale.GetTransactionByAddress(address, startBlock, endBlock)
	if err != nil {
		return
	}
	log.Println("start:", startBlock, "end:", endBlock, "len:", len(tx.Result), "err:", err)
	results = tx.Result

	addressTx := new(whale.AddressTransaction)
	addressTx.Address = address
	addressTx.Block = endBlock
	if len(results) > 0 {
		addressTx.TxId = results[len(results)-1].Hash
	}

	if startBlock == 0 {
		err = whale.DB().Create(addressTx).Error
		return
	}

	err = whale.DB().Model(addressTx).Where("address = ? ", address).UpdateColumn("block", endBlock).Error
	return
}

func ListenTx(cli *whale.Client) {
	for {
		time.Sleep(10 * time.Second)
		log.Println("start to listen....")
		endBlock, err := cli.BlockNumber()
		if err != nil {
			log.Fatal(err)
		}

		for _, addr := range addressList {
			results, err := scanTx(addr, endBlock)
			log.Println("address:", addr, "endBlock:", endBlock, "result len:", len(results))
			if err != nil {
				log.Println("addr:", addr, "err:", err)
				continue
			}
			for _, result := range results {
				err = whale.SendMsg(addr, result)
				log.Println("addr:", addr, "to:", to, "err:", err)
			}
		}
	}
}
