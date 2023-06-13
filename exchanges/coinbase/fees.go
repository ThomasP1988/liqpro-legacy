package coinbase

import "encoding/json"

type Fees struct {
	Maker     string `json:"maker_fee_rate"`
	Taker     string `json:"taker_fee_rate"`
	UsdVolume string `json:"usd_volume,omitempty"`
}

func GetFees() (*Fees, error) {
	body := []byte{}
	fees := &Fees{}
	response, errReponse := RequestSigned("GET", "/fees", &body)

	if errReponse != nil {
		return nil, errReponse
	}
	println("response", string(response))

	err := json.Unmarshal(response, fees)

	if err != nil {
		return nil, err
	}
	return fees, nil
}
