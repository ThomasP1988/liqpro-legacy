package kraken

import (
	"encoding/json"
	"net/url"
)

// GetFees get fees for desired pairs
func GetFees(symbols []string) (*TradeVolumeResponse, error) {

	// symbolsStr := strings.Join(symbols, ",")

	params := &url.Values{
		"pair": []string{"XBTUSD"},
	}

	response, errRequest := RequestSigned("TradeVolume", params)
	println("response", string(response))
	if errRequest != nil {
		println("kraken order request error", errRequest)
		return nil, errRequest
	}

	tradeVolume := &TradeVolumeResponse{}

	json.Unmarshal(response, tradeVolume)

	return tradeVolume, nil

}

func GetBalance() {
	params := &url.Values{}

	response, errRequest := RequestSigned("Balance", params)
	println("response", string(response))
	println("errRequest", errRequest)

}

// TradeVolumeResponse - response on TradeVolume request
type TradeVolumeResponse struct {
	Currency  string          `json:"currency"`
	Volume    float64         `json:"volume,string"`
	Fees      map[string]Fees `json:"fees,omitempty"`
	FeesMaker map[string]Fees `json:"fees_maker,omitempty"`
}

// Fees - structure of fees info
type Fees struct {
	Fee        float64 `json:"fee,string"`
	MinFee     float64 `json:"minfee,string"`
	MaxFee     float64 `json:"maxfee,string"`
	NextFee    float64 `json:"nextfee,string"`
	NextVolume float64 `json:"nextvolume,string"`
	TierVolume float64 `json:"tiervolume,string"`
}
