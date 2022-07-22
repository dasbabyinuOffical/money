package ERC20

import (
	"math/big"
)

type ERC20Handler interface {
	GetTotalSupply() (*big.Int, error)
	GetSymbol() string
	GetDecimal() uint8
	GetGasLimit(from, to string, value *big.Int) uint64
	GetBalance(wallet string) (*big.Int, error)
	GetBalanceByNumber(addr string, number int64) (*big.Int, error)
	Str2Big(amount string) (*big.Int, error)
	Big2Str(amount *big.Int) string
	Allowance(owner string, spender string) (amount *big.Int, err error)
}
