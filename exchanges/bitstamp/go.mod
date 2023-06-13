module liqpro/exchanges/bitstamp

go 1.16

require (
	github.com/andybalholm/brotli v1.0.1 // indirect
	github.com/fasthttp/websocket v1.4.3
	github.com/google/uuid v1.2.0 // indirect
	github.com/klauspost/compress v1.11.6 // indirect
	github.com/valyala/fasthttp v1.20.0 // indirect
	liqpro/exchanges/common v1.0.0
	liqpro/shared/libs/crypto v1.0.0
)

replace liqpro/shared/libs/crypto => ./../../shared/libs/crypto

replace liqpro/exchanges/common => ./../common
