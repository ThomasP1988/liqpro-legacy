module liqpro/pods/kraken

go 1.16

require (
	liqpro/exchanges/common v0.0.0-00010101000000-000000000000 // indirect
	liqpro/exchanges/kraken v1.0.0
	liqpro/shared/disruptor/parser v1.0.0
)

replace liqpro/exchanges/kraken => ./../../../exchanges/kraken

replace liqpro/shared/disruptor/parser => ./../../../shared/disruptor/parser

replace liqpro/exchanges/common => ./../../../exchanges/common

replace liqpro/shared/libs/crypto => ./../../../shared/libs/crypto
