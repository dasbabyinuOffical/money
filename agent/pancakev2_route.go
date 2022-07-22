package agent

import (
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
	"math/big"
	"money/agent/contracts"
	"strings"
)

const (
	swapExactETHForTokens = "swapExactETHForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline)"
)

type InputTx struct {
	// 卖出合约
	MakerContract string
	// 买入合约
	TakerContract string
}
type Decoder interface {
	DecodeInput(input string, logData []byte) (inputTx *InputTx, err error)
}

type SwapExactETHForTokensDecoder struct {
	Decoder
	FunctionInputs struct {
		AmountOutMin *big.Int
		Path         []common.Address
		To           common.Address
		Deadline     *big.Int
	}
}

func NewBscInputDecoder(funcName string) (decoder Decoder, err error) {
	switch funcName {
	case swapExactETHForTokens:
		decoder = &SwapExactETHForTokensDecoder{}
	default:
		err = errors.New("unsupport input decoder")
	}
	return
}

func (decoder *SwapExactETHForTokensDecoder) DecodeInput(input string, logData []byte) (inputTx *InputTx, err error) {
	inputTx = new(InputTx)
	abiContent := contracts.GetPancakeV2ABI()
	abi, err := abi.JSON(strings.NewReader(abiContent))
	if err != nil {
		return
	}

	decodedSig, err := hex.DecodeString(input[2:10])
	method, err := abi.MethodById(decodedSig)
	if err != nil {
		return
	}

	inputMap := make(map[string]interface{}, 0)
	data, err := hex.DecodeString(input[10:])
	if err != nil {
		return
	}

	err = method.Inputs.UnpackIntoMap(inputMap, data)
	if err != nil {
		return
	}

	err = mapstructure.Decode(inputMap, &decoder.FunctionInputs)
	if err != nil {
		return
	}
	path := decoder.FunctionInputs.Path
	if len(path) < 2 {
		err = errors.New("input path length error")
		return
	}

	inputTx.MakerContract = path[0].Hex()
	inputTx.TakerContract = path[len(path)-1].Hex()
	return
}
