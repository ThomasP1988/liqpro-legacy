module liqpro/pods/bitstamp

go 1.16

require (
	liqpro/exchanges/bitstamp v1.0.0
	liqpro/shared/disruptor/parser v1.0.0
	github.com/go-redis/redis/v8 v8.4.8 
)

replace liqpro/exchanges/bitstamp => ./../../../exchanges/bitstamp
replace liqpro/exchanges/common => ./../../../exchanges/common

replace liqpro/shared/disruptor/parser => ./../../../shared/disruptor/parser
replace liqpro/shared/libs/crypto => ./../../../shared/libs/crypto
