package main

import (
	"money/agent"
	"money/server"
	"time"
)

func Loop() {
	bscAgent, err := agent.NewAgent()
	if err != nil {
		return
	}
	for {
		time.Sleep(time.Second)
		bscAgent.ScanLog()
	}
}

func main() {
	go Loop()
	server.Serve()
}
