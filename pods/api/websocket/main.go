package main

import (
	"flag"
	"fmt"
	config "liqpro/config"
	"liqpro/shared/disruptor/order"
	"liqpro/shared/libs/auth"
	"liqpro/shared/repositories"
	"log"

	"github.com/valyala/fasthttp"
)

var addr = flag.String("addr", ":8080", "http service address")

// TheHub instantiation of the hub struct
var TheHub = NewHub()

func main() {
	config.SetConfig()
	ConnectToDB()
	repositories.ConnectRedisPrices()

	go func() {
		order.OrderDisruptor.Read()
	}()

	go (*TheHub).run()
	go myDisruptor.Read()

	requestHandler := func(ctx *fasthttp.RequestCtx) {

		// Authentication

		userDataCache, err := auth.GetUserIDFromQuery(ctx)

		if err != nil {
			fmt.Println("error", err)
			ctx.Response.SetStatusCode(1000)
			ctx.Write([]byte(err.Error()))
			return
		}

		// End Authentication

		switch string(ctx.Path()) {
		case "/ws":
			fmt.Println("/ws")
			HandleNewConnection(ctx, TheHub, userDataCache)
		default:
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}
	}

	server := fasthttp.Server{
		Name:    "Price",
		Handler: requestHandler,
	}

	log.Fatal(server.ListenAndServe(*addr))
}
