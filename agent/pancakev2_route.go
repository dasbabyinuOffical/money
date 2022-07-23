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
	swapExactETHForTokens                                 = "swapExactETHForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline)"
	swapExactTokensForETH                                 = "swapExactTokensForETH(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)"
	swapExactTokensForTokens                              = "swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)"
	swapExactTokensForTokensSupportingFeeOnTransferTokens = "swapExactTokensForTokensSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline)"
	addLiquidity                                          = "addLiquidity(address tokenA, address tokenB, uint256 amountADesired, uint256 amountBDesired, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline)"
)

type InputTx struct {
	// 卖出合约
	MakerContract string
	// 买入合约
	TakerContract string
}
type Decoder interface {
	DecodeInput(input string) (inputTx *InputTx, err error)
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

type SwapExactTokensForETHDecoder struct {
	Decoder
	FunctionInputs struct {
		AmountIn     *big.Int
		AmountOutMin *big.Int
		Path         []common.Address
		To           common.Address
		Deadline     *big.Int
	}
}

type SwapExactTokensForTokensDecoder struct {
	Decoder
	FunctionInputs struct {
		AmountIn     *big.Int
		AmountOutMin *big.Int
		Path         []common.Address
		To           common.Address
		Deadline     *big.Int
	}
}

type SwapExactTokensForTokensSupportingFeeOnTransferTokensDecoder struct {
	Decoder
	FunctionInputs struct {
		AmountIn     *big.Int
		AmountOutMin *big.Int
		Path         []common.Address
		To           common.Address
		Deadline     *big.Int
	}
}

type AddLiquidityDecoder struct {
	Decoder
	FunctionInputs struct {
		TokenA         common.Address
		TokenB         common.Address
		AmountADesired *big.Int
		AmountBDesired *big.Int
		AmountAMin     *big.Int
		AmountBMin     *big.Int
		To             common.Address
		Deadline       *big.Int
	}
}

func NewBscInputDecoder(funcName string) (decoder Decoder, err error) {
	switch funcName {
	case swapExactETHForTokens:
		decoder = new(SwapExactETHForTokensDecoder)
	case swapExactTokensForETH:
		decoder = new(SwapExactTokensForETHDecoder)
	case swapExactTokensForTokens:
		decoder = new(SwapExactTokensForTokensDecoder)
	case swapExactTokensForTokensSupportingFeeOnTransferTokens:
		decoder = new(SwapExactTokensForTokensSupportingFeeOnTransferTokensDecoder)
	case addLiquidity:
		decoder = new(AddLiquidityDecoder)
	default:
		err = errors.New("unsupport input decoder")
	}
	return
}

func (decoder *SwapExactETHForTokensDecoder) DecodeInput(input string) (inputTx *InputTx, err error) {
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

func (decoder *SwapExactTokensForETHDecoder) DecodeInput(input string) (inputTx *InputTx, err error) {
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

func (decoder *SwapExactTokensForTokensDecoder) DecodeInput(input string) (inputTx *InputTx, err error) {
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

func (decoder *SwapExactTokensForTokensSupportingFeeOnTransferTokensDecoder) DecodeInput(input string) (inputTx *InputTx, err error) {
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

func (decoder *AddLiquidityDecoder) DecodeInput(input string) (inputTx *InputTx, err error) {
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

	inputTx.MakerContract = decoder.FunctionInputs.TokenA.Hex()
	inputTx.TakerContract = decoder.FunctionInputs.TokenB.Hex()
	return
}
