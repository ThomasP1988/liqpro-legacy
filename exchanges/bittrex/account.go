package bittrex

import (
	"encoding/json"
	"strconv"
)

type AccountVolume struct {
	Updated      string `json:"updated"`
	Volume30days string `json:"volume30days"`
}

type FeeRate struct {
	Taker int
	Maker int
}

var FeesRate map[int]FeeRate = map[int]FeeRate{
	25000: { // 25k
		Taker: 35,
		Maker: 35,
	},
	50000: { // 50k
		Taker: 20,
		Maker: 25,
	},
	1000000: { // 1M
		Taker: 12,
		Maker: 18,
	},
	10000000: { // 10M
		Taker: 5,
		Maker: 15,
	},
	60000000: { // 60M
		Taker: 2,
		Maker: 10,
	},
	100000000: { // 100M
		Taker: 0,
		Maker: 8,
	},
	100000000000000: { // 100M+
		Taker: 0,
		Maker: 5,
	},
}

func GetFees() (*FeeRate, error) {

	acc, err := GetAccountVolume()

	if err != nil {
		println("error:", err)
		return nil, err
	}
	var volumeFlt float64
	volumeFlt, err = strconv.ParseFloat(acc.Volume30days, 64)

	if err != nil {
		println("error:", err)
		return nil, err
	}

	result := CalculateRate(volumeFlt, FeesRate)

	return result, nil

}

// GetAccountVolume get volume on account.
func GetAccountVolume() (*AccountVolume, error) {
	body, err := Request("GET", "/account/volume", nil, nil)
	if err != nil {
		println("error:", err)
		return nil, err
	}
	println("result:", string(*body))
	accVolume := &AccountVolume{}
	err = json.Unmarshal(*body, accVolume)
	if err != nil {
		println("error:", err)
		return nil, err
	}

	return accVolume, nil
}

func CalculateRate(volume float64, feesRates map[int]FeeRate) *FeeRate {

	for volumeKey, feeRate := range feesRates {
		if float64(volumeKey) > volume {
			println("FEE RATESSSSSS")
			return &feeRate
		}
	}

	return nil
}
