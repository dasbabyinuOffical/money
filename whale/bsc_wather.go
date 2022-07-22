package whale

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

const (
	PancakeV2RouterAddress = "0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d"
)

func WatchContractEvent() {
	client, err := ethclient.Dial("wss://bsc-mainnet.nodereal.io/ws/v1/2094fe91eea64b2dadcdced50fae154c")
	if err != nil {
		log.Fatal(err)
	}

	contractAddress := common.HexToAddress(PancakeV2RouterAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log
		}
	}
}
