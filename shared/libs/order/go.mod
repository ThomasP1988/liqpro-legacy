module liqpro/shared/libs/order

go 1.16

require (
	github.com/go-redis/redis/v8 v8.6.0
	github.com/google/uuid v1.2.0
	liqpro/exchanges/common v1.0.0
	liqpro/shared/disruptor/order v1.0.0
	liqpro/shared/repositories v0.0.0-00010101000000-000000000000
)

replace liqpro/shared/disruptor/order => ./../../disruptor/order

replace liqpro/exchanges/common => ./../../../exchanges/common

replace liqpro/exchanges/binance => ./../../../exchanges/binance

replace liqpro/exchanges/bitstamp => ./../../../exchanges/bitstamp

replace liqpro/exchanges/bittrex => ./../../../exchanges/bittrex

replace liqpro/exchanges/coinbase => ./../../../exchanges/coinbase

replace liqpro/exchanges/huobi => ./../../../exchanges/huobi

replace liqpro/shared/repositories => ./../../repositories

replace liqpro/shared/libs/crypto => ./../../libs/crypto

replace liqpro/shared/libs/signalr => ./../../libs/signalr
