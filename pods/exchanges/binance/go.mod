module liqpro/pods/binance

go 1.16

require (
	github.com/go-redis/redis/v8 v8.4.10
	liqpro/exchanges/binance v1.0.0
	liqpro/shared/disruptor/parser v1.0.0
)

replace liqpro/exchanges/binance => ./../../../exchanges/binance

replace liqpro/shared/disruptor/parser => ./../../../shared/disruptor/parser
replace liqpro/exchanges/common => ./../../../exchanges/common
