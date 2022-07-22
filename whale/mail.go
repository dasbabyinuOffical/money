package whale

import (
	"encoding/json"
)

const (
	TXURL   = "https://etherscan.io/tx/"
	ADDRURL = "https://etherscan.io/address/"
	AVEURL  = "https://ave.ai/token/"
	BOTURL  = "https://hooks.slack.com/services/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"
)

type BotMsg struct {
	Text string `json:"text"`
}

func SendMsg(addr string, result Result) (err error) {
	msg := BotMsg{Text: TXURL + result.Hash}
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	_, err = HttpPost(BOTURL, data)
	return
}
