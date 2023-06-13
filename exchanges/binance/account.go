package binance

import (
	"encoding/json"
	"fmt"
	"time"
)

func GetAccount() (*Account, error) {
	params := &map[string]string{
		"timestamp": fmt.Sprint(time.Now().Unix() * 1000),
	}
	response, err := RequestSigned("GET", "/account", params, "")
	println("err", err)
	println("response", string(response))
	if err != nil {
		return nil, err
	}
	account := &Account{}
	errJson := json.Unmarshal(response, account)

	if errJson != nil {
		return nil, errJson
	}

	return account, nil
}
