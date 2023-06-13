module liqpro/shared/repositories

go 1.16

require (
	github.com/dgraph-io/ristretto v0.0.3
	github.com/go-redis/redis/v8 v8.7.1
	go.mongodb.org/mongo-driver v1.4.5
	liqpro/config v0.0.0-00010101000000-000000000000
	liqpro/exchanges/binance v1.0.0
	liqpro/exchanges/bitstamp v1.0.0
	liqpro/exchanges/bittrex v1.0.0
	liqpro/exchanges/coinbase v1.0.0
	liqpro/exchanges/common v1.0.0
	liqpro/exchanges/huobi v1.0.0
	liqpro/shared/libs/crypto v1.0.0
)

replace liqpro/config => ./../../config

replace liqpro/shared/libs/signalr => ./../libs/signalr

replace liqpro/exchanges/binance => ./../../exchanges/binance

replace liqpro/exchanges/bitstamp => ./../../exchanges/bitstamp

replace liqpro/exchanges/bittrex => ./../../exchanges/bittrex

replace liqpro/exchanges/coinbase => ./../../exchanges/coinbase

replace liqpro/exchanges/huobi => ./../../exchanges/huobi

replace liqpro/exchanges/common => ./../../exchanges/common

replace liqpro/shared/libs/crypto => ./../libs/crypto
