// /v2/reference/transact-fee-rate

package huobi

import (
	"encoding/json"
	"net/url"
	"strings"
)

func GetFees(symbols []string) (*TransacFeeReponse, error) {

	symbolsStr := strings.Join(symbols, ",")

	params := &url.Values{
		"symbols": {
			symbolsStr,
		},
	}

	response, errRequest := RequestSigned("GET", "/v2/reference/transact-fee-rate", nil, params)
	println("response", string(response))
	if errRequest != nil {
		println("huobi order request error", errRequest)
		return nil, errRequest
	}

	transacFeeResponse := &TransacFeeReponse{}

	err := json.Unmarshal(response, transacFeeResponse)

	if err != nil {
		println("huobi order request error", err)
		return nil, err
	}

	return transacFeeResponse, nil
}

type TransacFeeReponse struct {
	Code int                     `json:"code"`
	Data []TransacFeeDataReponse `json:"data"`
}

type TransacFeeDataReponse struct {
	Symbol          string `json:"symbol"`
	ActualMakerRate string `json:"actualMakerRate"`
	ActualTakerRate string `json:"actualTakerRate"`
	TakerFeeRate    string `json:"takerFeeRate"`
	MakerFeeRate    string `json:"makerFeeRate"`
}
