module liqpro/exchanges/bittrex

go 1.16

require (
	github.com/andybalholm/brotli v1.0.1 // indirect
	github.com/klauspost/compress v1.11.6 // indirect
	github.com/valyala/fasthttp v1.20.0
	liqpro/shared/libs/signalr v1.0.0
		liqpro/exchanges/common v1.0.0
			liqpro/shared/libs/crypto v1.0.0

)

replace liqpro/shared/libs/signalr => ./../../shared/libs/signalr
replace liqpro/exchanges/common => ./../common
replace liqpro/shared/libs/crypto => ./../../shared/libs/crypto
