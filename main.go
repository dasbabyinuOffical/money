package main

import (
	"fmt"
	"money/agent"
	"money/server"
	"time"
)

func Loop(bscAgent *agent.Agent) {
	for {
		time.Sleep(time.Second)
		bscAgent.ScanLog()
	}
}

func SyncContractTask(bscAgent *agent.Agent) {
	for {
		time.Sleep(100 * time.Millisecond)
		err := bscAgent.SyncContractInfo()
		if err != nil {
			fmt.Println("SyncContractInfo,err:", err)
		}
	}
}

func main() {
	bscAgent, err := agent.NewAgent()
	if err != nil {
		panic(err)
	}
	go Loop(bscAgent)
	go SyncContractTask(bscAgent)
	server.Serve()
}
