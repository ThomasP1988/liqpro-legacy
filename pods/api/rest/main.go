package main

import (
	"bytes"
	"fmt"
	handlers "liqpro/pods/api/rest/handlers"
	"strings"

	config "liqpro/config"

	order "liqpro/shared/disruptor/order"

	fasthttp "github.com/valyala/fasthttp"

	repositories "liqpro/shared/repositories"

	cache "liqpro/shared/repositories/cache"

	"liqpro/shared/libs/auth"
)

var (
	headerNonce                     []byte = []byte("Nonce")
	headerSignature                 []byte = []byte("Signature")
	headerAPIKey                    []byte = []byte("API-Key")
	headerAccessControlAllowOrigin  []byte = []byte("Access-Control-Allow-Origin")
	headerAccessControlAllowHeaders []byte = []byte("Access-Control-Allow-Headers")
	headerOrigin                    []byte = []byte("Origin")
	pong                            []byte = []byte("pong")
	options                         []byte = []byte("OPTIONS")
	headersAuthorised               []byte = []byte(strings.Join([]string{string(headerAPIKey), string(headerNonce), string(headerSignature)}, ", "))
)

//HandleFastHTTP fastHTTP routing
func HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetBytesKV(headerAccessControlAllowHeaders, headersAuthorised)
	ctx.Response.Header.SetBytesKV(headerAccessControlAllowOrigin, ctx.Request.Header.PeekBytes(headerOrigin))

	if bytes.Compare(ctx.Request.Header.Method(), options) == 0 {
		ctx.Response.SetStatusCode(200)
		return

	}

	fmt.Println("request", string(ctx.Path()))
	userData, err := auth.GetUserDataWithHeadersFromCache(ctx)

	if err != nil {
		// return error
		fmt.Println("error", err)
		ctx.Response.SetStatusCode(400)
		ctx.Write([]byte(err.Error()))
		return
	}

	fmt.Println("userData", userData)
	switch string(ctx.Path()) {
	case "/instruments":
		handlers.ListInstruments()
	case "/order":
		handlers.Order(ctx, userData)
	case "/portfolio":
		handlers.Portfolio()
	case "/withdrawal":
		handlers.Withdrawal()
	case "/ping":
		ctx.Response.SetStatusCode(200)
		ctx.Response.SetBody(pong)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func main() {
	repositories.ConnectRedisPrices()
	config.SetConfig()

	errCache := cache.InitAPIUsers()

	if errCache != nil {
		panic(errCache)
	}

	// end cache

	go func() {
		order.OrderDisruptor.Read()
	}()

	fasthttp.ListenAndServe(":8083", HandleFastHTTP)
}
