module liqpro/pods/huobi

go 1.16

require (
	github.com/fasthttp/websocket v1.4.3
	github.com/go-redis/redis/v8 v8.4.10
	liqpro/exchanges/common v0.0.0-00010101000000-000000000000 // indirect
	liqpro/exchanges/huobi v1.0.0
	liqpro/shared/disruptor/parser v1.0.0
)

replace liqpro/exchanges/common => ./../../../exchanges/common

replace liqpro/exchanges/huobi => ./../../../exchanges/huobi

replace liqpro/shared/disruptor/parser => ./../../../shared/disruptor/parser
