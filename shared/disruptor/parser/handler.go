package socketredis

import (
	"context"
	"fmt"

	socket "github.com/fasthttp/websocket"
	redis "github.com/go-redis/redis/v8"
	disruptor "github.com/smartystreets-prototypes/go-disruptor"
)

// Define constant for the disruptor
const (
	BufferSize   = 64
	BufferMask   = BufferSize - 1
	Reservations = 1
)

// RingBuffer where we store the message before they are consumed
var RingBuffer = [BufferSize]*[]byte{}
var MyDisruptor disruptor.Disruptor

var (
	dbClient = redis.NewClient(&redis.Options{
		Addr:     "192.168.49.2:31928",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	ctx      = context.Background()
	Pipeline redis.Pipeliner
)

// DisruptorBytes same as disruptor, but the event handler dont handle the socket, only receive bytes
func DisruptorBytes(consumers ...disruptor.Consumer) WsHandlerBytes {

	MyDisruptor = disruptor.New(
		disruptor.WithCapacity(BufferSize),
		disruptor.WithConsumerGroup(consumers...))

	Pipeline = dbClient.Pipeline()
	// pipeline.Set(ctx, "test", "", 60*60)
	go func() {
		MyDisruptor.Read()
	}()

	return WsDepthHandlerBytes
	// _ = MyDisruptor.Close()
}

// sequence if the loop is handled outside, we need to have a sequence here
var sequence int64 = 0

// WsDepthHandlerBytes taking the msg from WS and sending it to our threads
func WsDepthHandlerBytes(msg *[]byte) {
	sequence = MyDisruptor.Reserve(Reservations)
	RingBuffer[sequence&BufferMask] = msg
	MyDisruptor.Commit(sequence-Reservations+1, sequence)
}

// Disruptor pods should call this function
func Disruptor(consumers ...disruptor.Consumer) WsHandler {

	MyDisruptor = disruptor.New(
		disruptor.WithWaitStrategy(NewFastWaitStrategy()),
		disruptor.WithCapacity(BufferSize),
		disruptor.WithConsumerGroup(consumers...))

	Pipeline = dbClient.Pipeline()
	// pipeline.Set(ctx, "test", "", 60*60)
	go func() {
		MyDisruptor.Read()
	}()

	return WsDepthHandler
	// _ = MyDisruptor.Close()
}

// WsDepthHandler taking the msg from WS and sending it to our threads
func WsDepthHandler(c *socket.Conn) {

	var sequence int64 = 0
	var messages = make([][]byte, 60)
	var i = 0
	var err error

	for {

		sequence = MyDisruptor.Reserve(Reservations)
		_, messages[i], err = c.ReadMessage()
		// println("message", string(messages[i]))
		if err != nil {
			fmt.Println("error:", err)
			break
		}
		RingBuffer[sequence&BufferMask] = &messages[i]
		MyDisruptor.Commit(sequence-Reservations+1, sequence)

		i++
		if i > 58 {
			i = 0
		}
	}
}

// Upsert data to DB
func Upsert(b *ConsumerBase) {
	b.Prefix = b.Market + ":" + b.Action

	println("b.Prefix", b.Prefix)
	println("b.Score", b.Score)

	Pipeline.Process(ctx, redis.NewIntCmd(ctx, "zadd", b.Prefix, b.Score, b.Volume+":"+b.Price+":"+b.Platform))

	b.RedisResult, b.Err = Pipeline.Exec(ctx)

	if b.Err != nil {
		fmt.Println("error upsert ", b.Err)
	} else {
		fmt.Print("+")
	}

}

// Del data to DB
func Del(b *ConsumerBase) {
	b.Prefix = b.Market + ":" + b.Action
	b.Key = b.Volume + ":" + b.Price + ":" + b.Platform

	Pipeline.ZRem(ctx, b.Prefix, b.Key)

	b.RedisResult, b.Err = Pipeline.Exec(ctx)
	if b.Err != nil {
		fmt.Println("error", b.Err)
	} else {
		fmt.Print("-")
	}

}

func StopDisruptor() {
	MyDisruptor.Close()
}
