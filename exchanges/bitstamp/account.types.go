package bitstamp

type ErrorResult struct {
	Status string `json:"status,string"`
	Reason string `json:"reason,string"`
	Code   string `json:"code,string"`
}

type AccountBalanceResult struct {
	UsdBalance   float64 `json:"usd_balance,string"`
	BtcBalance   float64 `json:"btc_balance,string"`
	EurBalance   float64 `json:"eur_balance,string"`
	XrpBalance   float64 `json:"xrp_balance,string"`
	LtcBalance   float64 `json:"ltc_balance,string"`
	EthBalance   float64 `json:"eth_balance,string"`
	BchBalance   float64 `json:"bch_balance,string"`
	UsdReserved  float64 `json:"usd_reserved,string"`
	BtcReserved  float64 `json:"btc_reserved,string"`
	EurReserved  float64 `json:"eur_reserved,string"`
	XrpReserved  float64 `json:"xrp_reserved,string"`
	LtcReserved  float64 `json:"ltc_reserved,string"`
	EthReserved  float64 `json:"eth_reserved,string"`
	BchReserved  float64 `json:"bch_reserved,string"`
	UsdAvailable float64 `json:"usd_available,string"`
	BtcAvailable float64 `json:"btc_available,string"`
	EurAvailable float64 `json:"eur_available,string"`
	XrpAvailable float64 `json:"xrp_available,string"`
	LtcAvailable float64 `json:"ltc_available,string"`
	EthAvailable float64 `json:"eth_available,string"`
	BchAvailable float64 `json:"bch_available,string"`
	BtcUsdFee    float64 `json:"btcusd_fee,string"`
	BtcEurFee    float64 `json:"btceur_fee,string"`
	EurUsdFee    float64 `json:"eurusd_fee,string"`
	XrpUsdFee    float64 `json:"xrpusd_fee,string"`
	XrpEurFee    float64 `json:"xrpeur_fee,string"`
	XrpBtcFee    float64 `json:"xrpbtc_fee,string"`
	LtcUsdFee    float64 `json:"ltcusd_fee,string"`
	LtcEurFee    float64 `json:"ltceur_fee,string"`
	LtcBtcFee    float64 `json:"ltcbtc_fee,string"`
	EthUsdFee    float64 `json:"ethusd_fee,string"`
	EthEurFee    float64 `json:"etheur_fee,string"`
	EthBtcFee    float64 `json:"ethbtc_fee,string"`
	BchUsdFee    float64 `json:"bchusd_fee,string"`
	BchEurFee    float64 `json:"bcheur_fee,string"`
	BchBtcFee    float64 `json:"bchbtc_fee,string"`
}
