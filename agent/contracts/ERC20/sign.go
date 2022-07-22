package ERC20

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/inwecrypto/sha3"
)

const (
	signBalanceOf    = "balanceOf(address)"
	signTotalSupply  = "totalSupply()"
	signTransfer     = "transfer(address,uint256)"
	signTransferFrom = "transferFrom(address,address,uint256)"
	signApprove      = "approve(address,uint256)"
	signName         = "name()"
	signSymbol       = "symbol()"
	signAllowance    = "allowance(address,address)"
	signdecimals     = "decimals()"
)

// Method/Event id
var (
	TransferID     = SignABI(signTransfer)
	BalanceOfID    = SignABI(signBalanceOf)
	Decimals       = SignABI(signdecimals)
	TransferFromID = SignABI(signTransferFrom)
	ApproveID      = SignABI(signApprove)
	TotalSupplyID  = SignABI(signTotalSupply)
	AllowanceID    = SignABI(signAllowance)
)

// SignABI sign abi string
func SignABI(abi string) string {
	hasher := sha3.NewKeccak256()
	hasher.Write([]byte(abi))
	data := hasher.Sum(nil)

	return hex.EncodeToString(data[0:4])
}

func packNumeric(value string, bytes int) string {
	if value == "" {
		value = "0x0"
	}

	value = strings.TrimPrefix(value, "0x")

	chars := bytes * 2

	n := len(value)
	if n%chars == 0 {
		return value
	}
	return strings.Repeat("0", chars-n%chars) + value
}

// BalanceOf create erc20 balanceof abi string
func BalanceOf(address string) string {
	address = packNumeric(address, 32)

	return fmt.Sprintf("%s%s", BalanceOfID, address)
}

// Transfer .
func TransferHex(to string, value string) ([]byte, error) {
	to = packNumeric(to, 32)
	value = packNumeric(value, 32)

	data := fmt.Sprintf("%s%s%s", SignABI(signTransfer), to, value)

	return hex.DecodeString(data)
}
