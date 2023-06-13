package handlers

import (
	"encoding/json"
	"fmt"

	"liqpro/shared/libs/order"
	"liqpro/shared/repositories/cache"

	"github.com/valyala/fasthttp"
)

// Order clients pass orderRequest here
func Order(ctx *fasthttp.RequestCtx, userData *cache.UserDataCache) {
	orderRequest := &order.RequestArgs{}
	err := json.Unmarshal(ctx.PostBody(), orderRequest)
	fmt.Println("orderRequest", orderRequest)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	bson, err := order.Handle(orderRequest, userData)

	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.Response.SetBody([]byte(err.Error()))
		return
	}

	ctx.Response.SetStatusCode(200)
	ctx.Write(bson)
	ctx.SetContentType("application/json")
}
