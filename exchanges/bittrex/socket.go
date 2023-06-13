package bittrex

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"liqpro/shared/libs/signalr"
)

var (
	baseWsURL = "socket-v3.bittrex.com"
)

// doAsyncTimeout runs f in a different goroutine
//	if f returns before timeout elapses, doAsyncTimeout returns the result of f().
//	otherwise it returns "operation timeout" error, and calls tmFunc after f returns.
func doAsyncTimeout(f func() error, tmFunc func(error), timeout time.Duration) error {
	errs := make(chan error)
	go func() {
		err := f()
		select {
		case errs <- err:
		default:
			if tmFunc != nil {
				tmFunc(err)
			}
		}
	}()
	select {
	case err := <-errs:
		return err
	case <-time.After(timeout):
		return errors.New("operation timeout")
	}
}

// WsDepthServe connect to bittrex stream
func WsDepthServe(symbols []string, handler WsDepthHandler) error {
	const timeout = 5 * time.Second
	client := signalr.NewWebsocketClient()
	// defer client.Close()

	client.OnClientMethod = func(hub, method string, messages []json.RawMessage) {
		for _, msg := range messages {
			// the base64 is quoted, so we need to remove first byte and last
			msgLen := len(msg)
			result := make([]byte, base64.StdEncoding.EncodedLen(msgLen-2))
			info, err := base64.StdEncoding.Decode(result, msg[1:msgLen-1])
			if err != nil {
				fmt.Println("info", info)
				fmt.Println("err", err)
				continue
			}
			deflated, err := ioutil.ReadAll(flate.NewReader(bytes.NewReader(result)))

			handler(&deflated)
		}
	}

	err := client.Connect("https", baseWsURL, []string{"c3"})
	if err != nil {
		fmt.Println("error connecting: ", err)
	}
	fmt.Println("Connected")
	channels := []string{}
	for _, symbol := range symbols {
		channels = append(channels, "orderbook_"+symbol+"_25")
	}
	msg, err2 := client.CallHub("c3", "subscribe", channels)
	if err2 != nil {
		fmt.Println("error calling hub", err2)
	}
	fmt.Println("subscribe", string(msg))
	return nil

	// err := doAsyncTimeout(func() error {
	// 	return client.Connect("https", "api.bittrex.com/v3", []string{"c3"})
	// }, func(err error) {
	// 	if err == nil {
	// 		client.Close()
	// 	}
	// }, timeout)

	// if err != nil {
	// 	return err
	// }
	// defer client.Close()

	// err = doAsyncTimeout(func() error {
	// 	msg, err := client.CallHub("c3", "orderbook_"+symbol+"_25")
	// 	fmt.Println(msg)
	// 	return err
	// }, nil, timeout)

	// if err != nil {
	// 	return err
	// }

	// return nil
}
