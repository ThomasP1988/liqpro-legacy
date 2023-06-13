package config

// // TradingPair enum of currencies
// type TradingPair string

// // list of trading pairs
// const (
// 	BTCEUR TradingPair = "btceur"
// 	BTCUSD TradingPair = "btcusd"
// )

// Currency enum of currencies
type Currency string

// list of currencyes
const (
	BTC Currency = "btc"
	ETH Currency = "eth"
	EUR Currency = "eur"
	USD Currency = "usd"
)

// PairInfo centralize all information about a pair
type PairInfo struct {
	Base         Currency // first
	Quote        Currency // second
	BaseDecimal  int
	QuoteDecinal int
}
