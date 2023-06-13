package repositories

import (
	binance "liqpro/exchanges/binance"
	bitstamp "liqpro/exchanges/bitstamp"
	bittrex "liqpro/exchanges/bittrex"
	coinbase "liqpro/exchanges/coinbase"
	commonExchanges "liqpro/exchanges/common"
	huobi "liqpro/exchanges/huobi"
)

// Platforms set all the platform adapter in this struct
type Platforms struct {
	Platforms map[string]*commonExchanges.Platform
}

// Get platform adapter
func (ps *Platforms) Get(platformName string) *commonExchanges.Platform {
	return ps.Platforms[platformName]
}

// PlatformFactory should deliver platform adapter on demand
var PlatformFactory *Platforms = &Platforms{
	Platforms: map[string]*commonExchanges.Platform{
		binance.Adapter.Name():  &binance.Adapter,
		bitstamp.Adapter.Name(): &bitstamp.Adapter,
		bittrex.Adapter.Name():  &bittrex.Adapter,
		coinbase.Adapter.Name(): &coinbase.Adapter,
		huobi.Adapter.Name():    &huobi.Adapter,
	},
}
