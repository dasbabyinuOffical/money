package tools

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math/big"
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Uuid(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Uuid32() string {
	return Uuid(32)
}

func CloneInt(x *big.Int) *big.Int {
	return new(big.Int).Set(x)
}

func Big2Str(x *big.Int, decimal int) string {
	if decimal == 0 {
		return x.String()
	}

	// Integral part
	i := x.String()
	if len(i) <= decimal {
		i = "0"
	} else {
		i = i[0 : len(i)-decimal]
	}

	// Decimal part
	d := x.String()
	if len(d) < decimal {
		d = strings.Repeat("0", decimal-len(d)) + d
	} else {
		d = d[len(d)-decimal:]
	}

	return i + "." + d
}

func MulDecimal(amount string, decimals int) string {
	if !IsNumeric(amount) {
		return "0"
	}

	index := strings.Index(amount, ".")
	decimalsLen := 0

	if index >= 0 {
		decimalsLen = len(amount[index+1:])
		amount = amount[:index] + amount[index+1:]
	}

	if decimals > decimalsLen {
		amount = amount + strings.Repeat("0", decimals-decimalsLen)
	} else {
		amount = amount[:len(amount)-(decimalsLen-decimals)]
	}

	amount = strings.TrimLeft(amount, "0")
	if amount == "" {
		amount = "0"
	}

	return amount
}

func DivDecimal(amount string, decimal int) string {
	if amount == "" || !IsNumeric(amount) {
		return "0"
	}

	x, _ := big.NewInt(0).SetString(amount, 10)

	if decimal == 0 {
		return x.String()
	}

	// Integral part
	i := x.String()
	if len(i) <= decimal {
		i = "0"
	} else {
		i = i[0 : len(i)-decimal]
	}

	// Decimal part
	d := x.String()
	if len(d) < decimal {
		d = strings.Repeat("0", decimal-len(d)) + d
	} else {
		d = d[len(d)-decimal:]
	}

	return i + "." + d
}

func Str2Big(amount string, decimals int) (*big.Int, error) {
	amt, _, err := big.ParseFloat(amount, 10, 256, big.ToNearestEven)
	if err != nil || amt.Sign() < 0 {
		return nil, errors.New("Failed to parse amount: " + amount)
	}

	index := strings.Index(amount, ".")
	decimalsLen := 0

	if index >= 0 {
		decimalsLen = len(amount[index+1:])
		amount = amount[:index] + amount[index+1:]
	}

	if decimals > decimalsLen {
		amount = amount + strings.Repeat("0", decimals-decimalsLen)
	} else {
		amount = amount[:len(amount)-(decimalsLen-decimals)]
	}

	res, _ := big.NewInt(0).SetString(string(amount), 10)
	floatRes, _ := big.NewFloat(0).SetString(string(amount))

	base, _ := big.NewFloat(0).SetString("1" + strings.Repeat("0", decimals))
	amt = amt.Mul(amt, base)

	max := big.NewFloat(0).Mul(amt, big.NewFloat(1.1))
	min := big.NewFloat(0).Mul(amt, big.NewFloat(0.9))

	if floatRes.Cmp(max) <= 0 && floatRes.Cmp(min) >= 0 {
		return res, nil
	} else {
		bigAmt, _ := amt.Int(nil)
		return bigAmt, nil
	}
}

func Ask4confirm(msg string) bool {
	var s string

	fmt.Println(msg)
	fmt.Printf("(Y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

func GetSha256Hash(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = strings.Trim(str, " \\t\\n\\r\\v\\f")
		if str[0] == '-' || str[0] == '+' {
			if len(str) == 1 {
				return false
			}
			str = str[1:]
		}
		// hex
		if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
			for _, h := range str[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}
		// 0-9,Point,Scientific
		p, s, l := 0, 0, len(str)
		for i, v := range str {
			if v == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v == 'e' || v == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}
				s = i
			} else if v < '0' || v > '9' {
				return false
			}
		}
		return true
	}

	return false
}

// ParseBigInt parse hex string value to big.Int
func ParseBigInt(value string) (*big.Int, error) {
	i := big.NewInt(0)
	_, err := fmt.Sscan(value, i)

	return i, err
}

func SignABI(abi string) string {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(abi))
	data := hasher.Sum(nil)

	return hex.EncodeToString(data[0:4])
}
