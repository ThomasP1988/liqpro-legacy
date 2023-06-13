module liqpro/pods/coinbase

go 1.16

require (
	github.com/go-redis/redis/v8 v8.6.0
	liqpro/exchanges/coinbase v1.0.0
	liqpro/exchanges/common v0.0.0-00010101000000-000000000000 // indirect
	liqpro/shared/disruptor/parser v1.0.0
)

replace liqpro/exchanges/coinbase => ./../../../exchanges/coinbase

replace liqpro/exchanges/common => ./../../../exchanges/common

replace liqpro/shared/disruptor/parser => ./../../../shared/disruptor/parser
