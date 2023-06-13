module liqpro/shared/disruptor/order

go 1.16

require (
	github.com/smartystreets-prototypes/go-disruptor v0.0.0-20200316140655-c96477fd7a6a // indirect
    liqpro/shared/repositories v0.0.0-00010101000000-000000000000
    liqpro/exchanges/common v1.0.0
    
)

replace liqpro/shared/repositories => ./../../repositories
replace liqpro/exchanges/common => ./../../../exchanges/common


