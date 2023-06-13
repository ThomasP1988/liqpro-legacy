package bitstamp

import (
	"encoding/json"
	"net/url"
)

func GetAccount() (*AccountBalanceResult, error) {
	balance := &AccountBalanceResult{}
	data, err := RequestSigned("POST", "/balance/", url.Values{})
	if err != nil {
		return nil, err
	}
	println("data", data)
	err = json.Unmarshal(data, balance)

	if err != nil {
		return nil, err
	}

	return balance, nil
}
