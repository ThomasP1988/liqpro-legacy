module liqpro/pods/bittrex

go 1.16

require (
	github.com/go-redis/redis/v8 v8.4.10 // indirect
	liqpro/exchanges/bittrex v1.0.0
	liqpro/shared/disruptor/parser v0.0.0-00010101000000-000000000000 // indirect
)

replace liqpro/exchanges/bittrex => ./../../../exchanges/bittrex

replace liqpro/shared/libs/signalr => ./../../../shared/libs/signalr

replace liqpro/shared/disruptor/parser => ./../../../shared/disruptor/parser

replace liqpro/exchanges/common => ./../../../exchanges/common

replace liqpro/shared/libs/crypto => ./../../../shared/libs/crypto

