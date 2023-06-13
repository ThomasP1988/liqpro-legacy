module liqpro/pods/api/rest

go 1.16

require (
	github.com/klauspost/compress v1.11.7 // indirect
	github.com/valyala/fasthttp v1.20.0
	liqpro/config v1.0.0
	liqpro/shared/disruptor/order v1.0.0
	liqpro/shared/libs/auth v1.0.0
	liqpro/shared/libs/order v1.0.0
	liqpro/shared/repositories v0.0.0-00010101000000-000000000000
)

replace liqpro/config => ./../../../config

replace liqpro/shared/disruptor/order => ./../../../shared/disruptor/order

replace liqpro/shared/libs/crypto => ./../../../shared/libs/crypto

replace liqpro/shared/libs/auth => ./../../../shared/libs/auth

replace liqpro/shared/libs/order => ./../../../shared/libs/order

replace liqpro/exchanges/common => ./../../../exchanges/common

replace liqpro/exchanges/binance => ./../../../exchanges/binance

replace liqpro/exchanges/bitstamp => ./../../../exchanges/bitstamp

replace liqpro/exchanges/bittrex => ./../../../exchanges/bittrex

replace liqpro/exchanges/coinbase => ./../../../exchanges/coinbase

replace liqpro/exchanges/huobi => ./../../../exchanges/huobi

replace liqpro/shared/repositories => ./../../../shared/repositories

replace liqpro/shared/libs/signalr => ./../../../shared/libs/signalr
