module liqpro/exchanges/binance

go 1.16

require (
	github.com/fasthttp/websocket v1.4.3
	github.com/valyala/fasthttp v1.20.0
	liqpro/exchanges/common v1.0.0
)

replace liqpro/exchanges/common => ./../common
