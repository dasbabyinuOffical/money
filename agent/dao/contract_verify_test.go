package dao

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnalysisContract(t *testing.T) {
	contract := "0x2c9ecE873Bb1aD75D43b236879b0b49DEA4AF072"
	verify := new(ContractVerifyResult)
	data := `{
    "code": 1,
    "message": "OK",
    "result": {
        "0x2c9ece873bb1ad75d43b236879b0b49dea4af072": {
            "buy_tax": "0",
            "can_take_back_ownership": "0",
            "cannot_buy": "0",
            "cannot_sell_all": "0",
            "creator_address": "0x863b49ae97c3d2a87fd43186dfd921f42783c853",
            "creator_balance": "0",
            "creator_percent": "0",
            "dex": [
                {
                    "name": "PancakeV2",
                    "liquidity": "1918.74839923",
                    "pair": "0x7e7B2b0B10FB918e9b3A8DeF13acBC5Dab727108"
                }
            ],
            "hidden_owner": "0",
            "holder_count": "128",
            "holders": [
                {
                    "address": "0xe1b51f15688839351f7ae216416470bb7c57c984",
                    "tag": "",
                    "is_contract": 0,
                    "balance": "350609637342.39",
                    "percent": "0.438262046677987500",
                    "is_locked": 0
                },
                {
                    "address": "0x407993575c91ce7643a4d4ccacc9a98c36ee1bbe",
                    "tag": "PinkLock02",
                    "is_contract": 1,
                    "balance": "160000000000",
                    "percent": "0.200000000000000000",
                    "is_locked": 1,
                    "locked_detail": [
                        {
                            "amount": "160000000000.000000000000000000",
                            "end_time": "2032-07-01T11:35:00+00:00",
                            "opt_time": "2022-07-25T11:35:40+00:00"
                        }
                    ]
                },
                {
                    "address": "0x7e7b2b0b10fb918e9b3a8def13acbc5dab727108",
                    "tag": "PancakeV2",
                    "is_contract": 1,
                    "balance": "109984017461.65",
                    "percent": "0.137480021827062500",
                    "is_locked": 0
                },
                {
                    "address": "0x9276cd47ba544a4247ed89eb77226dfea63f1e03",
                    "tag": "",
                    "is_contract": 0,
                    "balance": "16234364778.809",
                    "percent": "0.020292955973511250",
                    "is_locked": 0
                },
                {
                    "address": "0x6408d8aff0d91369360f7fb8bad5a3c63e329f56",
                    "tag": "",
                    "is_contract": 0,
                    "balance": "13821572050.849",
                    "percent": "0.017276965063561250",
                    "is_locked": 0
                },
                {
                    "address": "0xf128598da91b7ae26f4c5adc4d52ce047ff682d3",
                    "tag": "",
                    "is_contract": 0,
                    "balance": "8132928777.6894",
                    "percent": "0.010166160972111750",
                    "is_locked": 0
                },
                {
                    "address": "0x21f9c9d849d49f38f1c04cdc4b6a65ab12ae2306",
                    "tag": "",
                    "is_contract": 0,
                    "balance": "6804478947.1654",
                    "percent": "0.008505598683956750",
                    "is_locked": 0
                },
                {
                    "address": "0x786a4024b366e8eef16e12e05f970a65d2c5fd85",
                    "tag": "",
                    "is_contract": 0,
                    "balance": "5588377981.6881",
                    "percent": "0.006985472477110125",
                    "is_locked": 0
                },
                {
                    "address": "0x3353fab5bdf866d6fc3610baf61392088625091a",
                    "tag": "",
                    "is_contract": 0,
                    "balance": "4815079939.9572",
                    "percent": "0.006018849924946500",
                    "is_locked": 0
                },
                {
                    "address": "0xb2e1fddd1cba010d81e6fd78cafa2227349978fd",
                    "tag": "",
                    "is_contract": 0,
                    "balance": "4211656956.021",
                    "percent": "0.005264571195026250",
                    "is_locked": 0
                }
            ],
            "is_anti_whale": "0",
            "is_blacklisted": "0",
            "is_honeypot": "0",
            "is_in_dex": "1",
            "is_mintable": "0",
            "is_open_source": "1",
            "is_proxy": "0",
            "is_whitelisted": "0",
            "lp_holder_count": "2",
            "lp_holders": [
                {
                    "address": "0x407993575c91ce7643a4d4ccacc9a98c36ee1bbe",
                    "tag": "PinkLock02",
                    "is_contract": 1,
                    "balance": "14419986.130368",
                    "percent": "1.000000000000013193",
                    "is_locked": 1,
                    "locked_detail": [
                        {
                            "amount": "14419986.130367809753913954",
                            "end_time": "2022-08-24T13:11:34+00:00",
                            "opt_time": "2022-07-25T13:11:34+00:00"
                        }
                    ]
                },
                {
                    "address": "0x0000000000000000000000000000000000000000",
                    "tag": "",
                    "is_contract": 0,
                    "balance": "0.000000000000001000",
                    "percent": "0.000000000000000000",
                    "is_locked": 1
                }
            ],
            "lp_total_supply": "14419986.130367809753914954",
            "owner_address": "",
            "owner_balance": "0",
            "owner_change_balance": "0",
            "owner_percent": "0",
            "sell_tax": "0",
            "slippage_modifiable": "0",
            "total_supply": "800000000000.000000000000000000",
            "trading_cooldown": "0",
            "transfer_pausable": "0",
            "token_name": "Vajra kangaroo",
            "token_symbol": "VAKG"
        }
    }
}`
	err = json.Unmarshal([]byte(data), verify)
	assert.Nil(t, err)
	gotScore, err := AnalysisContract(contract, verify)
	assert.Nil(t, err)
	t.Log(gotScore)
}
