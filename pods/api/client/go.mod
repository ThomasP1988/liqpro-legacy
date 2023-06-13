module liqpro/pods/api/client

go 1.16

require (
	github.com/gofiber/adaptor/v2 v2.1.1
	github.com/gofiber/fiber/v2 v2.4.1
	github.com/graphql-go/graphql v0.7.9
	github.com/pkg/errors v0.9.1
	github.com/supertokens/supertokens-go v1.4.2
	liqpro/config v0.0.0-00010101000000-000000000000 // indirect
	liqpro/shared/libs/crypto v1.0.0
	liqpro/shared/repositories v0.0.0-00010101000000-000000000000
)

replace liqpro/config => ./../../../config

replace liqpro/shared/repositories => ./../../../shared/repositories

replace liqpro/shared/libs/crypto => ./../../../shared/libs/crypto

replace liqpro/exchanges/common => ./../../../exchanges/common

replace liqpro/exchanges/binance => ./../../../exchanges/binance

replace liqpro/exchanges/bitstamp => ./../../../exchanges/bitstamp

replace liqpro/exchanges/bittrex => ./../../../exchanges/bittrex

replace liqpro/exchanges/coinbase => ./../../../exchanges/coinbase

replace liqpro/exchanges/huobi => ./../../../exchanges/huobi

replace liqpro/shared/libs/signalr => ./../../../shared/libs/signalr
