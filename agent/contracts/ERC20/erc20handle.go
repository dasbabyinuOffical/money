package ERC20

import (
	"context"
	"encoding/hex"
	"math/big"
	"money/agent/tools"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type ERC20Handle struct {
	Symbol      string
	Contract    string
	Decimals    uint8
	TotalSupply *big.Int
	Client      *ethclient.Client
	Transactor  *ERC20Transactor
	Caller      *ERC20Caller
}

func NewHandle(url string, contract string) (*ERC20Handle, error) {
	rpcDial, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}

	client := ethclient.NewClient(rpcDial)

	transactor, err := NewERC20Transactor(common.HexToAddress(contract), client)
	if err != nil {
		return nil, err
	}

	caller, err := NewERC20Caller(common.HexToAddress(contract), client)
	if err != nil {
		return nil, err
	}

	symbol, err := caller.Symbol(nil)
	if err != nil {
		return nil, err
	}

	decimals, err := caller.Decimals(nil)
	if err != nil {
		return nil, err
	}

	totalSupply, err := caller.TotalSupply(nil)
	if err != nil {
		return nil, err
	}

	h := &ERC20Handle{
		Symbol:      strings.ToUpper(symbol),
		TotalSupply: totalSupply,
		Decimals:    decimals,
		Client:      client,
		Transactor:  transactor,
		Caller:      caller,
		Contract:    contract,
	}

	return h, nil
}

func (h *ERC20Handle) GetSymbol() string {
	return h.Symbol
}

func (h *ERC20Handle) GetDecimal() uint8 {
	return h.Decimals
}

func (h *ERC20Handle) GetTotalSupply() (*big.Int, error) {
	return h.Caller.TotalSupply(nil)
}

func (h *ERC20Handle) GetGasLimit(from, to string, value *big.Int) uint64 {
	contractAddr := common.HexToAddress(h.Contract)
	methodData, err := TransferHex(to, hex.EncodeToString(value.Bytes()))
	if err != nil {
		return 90000
	}

	msg := ethereum.CallMsg{
		From:  common.HexToAddress(from),
		To:    &contractAddr,
		Value: nil,
		Data:  methodData,
	}

	gaslimit, err := h.Client.EstimateGas(context.Background(), msg)
	if err != nil {
		return 90000
	}

	return gaslimit + 5000
}

func (h *ERC20Handle) GetBalance(address string) (*big.Int, error) {
	return h.Caller.BalanceOf(nil, common.HexToAddress(address))
}

func (h *ERC20Handle) Str2Big(amount string) (*big.Int, error) {
	return tools.Str2Big(amount, int(h.Decimals))
}

func (h *ERC20Handle) Big2Str(amount *big.Int) string {
	return tools.Big2Str(amount, int(h.Decimals))
}

func (h *ERC20Handle) Allowance(owner string, spender string) (amount *big.Int, err error) {
	amount, err = h.Caller.Allowance(nil, common.HexToAddress(owner), common.HexToAddress(spender))
	return
}

func (h *ERC20Handle) GetBalanceByNumber(address string, num int64) (*big.Int, error) {

	contractAddr := common.HexToAddress(h.Contract)

	methodData := BalanceOf(address)

	data, err := hex.DecodeString(methodData)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		From: common.HexToAddress(address),
		To:   &contractAddr,
		Data: data,
	}

	bytes, err := h.Client.CallContract(context.TODO(), msg, big.NewInt(num))
	if nil != err {
		bytes, err = h.Client.CallContract(context.TODO(), msg, nil)
		if err != nil {
			return nil, err
		}
	}

	balance := big.NewInt(0).SetBytes(bytes)

	return balance, nil
}
