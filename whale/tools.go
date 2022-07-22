package whale

import (
	"bytes"
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"net/http"
	"time"
)

func HttpGet(url string) (data []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		return
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	data, err = ioutil.ReadAll(resp.Body)
	return
}

func HttpPost(url string, buf []byte) (data []byte, err error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		return
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	data, err = ioutil.ReadAll(resp.Body)
	return
}

type Client struct {
	client *ethclient.Client
}

const (
	EthUrl = "https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
)

func NewClient() (client *Client, err error) {
	ethClient, err := ethclient.Dial(EthUrl)
	if err != nil {
		return
	}
	client = &Client{
		client: ethClient,
	}
	return
}

func (cli *Client) Close() {
	if cli.client != nil {
		cli.client.Close()
	}
}

func (cli *Client) BlockNumber() (num uint64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	num, err = cli.client.BlockNumber(ctx)
	return
}
