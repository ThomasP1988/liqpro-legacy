package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/valyala/fasthttp"
)

const (
	baseURL = "https://api.binance.com/api/v3"
)

// APIKey Binance API credentials
const APIKey string = ""

// SecretKey Binance API credentials
const SecretKey string = ""

// APIKeyBytes Binance API credentials
var APIKeyBytes = []byte(APIKey)

// SecretKeyBytes Binance API credentials
var SecretKeyBytes = []byte(SecretKey)

// Request send REST request to bittrex
func Request(method string, endpoint string, params map[string]string, content string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	for key, val := range params {
		req.URI().QueryArgs().Add(key, val)
	}

	contentByte := []byte(content)

	req.Header.SetMethod(method)
	req.SetRequestURI(baseURL + endpoint)
	req.SetBody(contentByte)

	err := fasthttp.Do(req, resp)
	bodyBytes := resp.Body()
	return bodyBytes, err
}

// RequestSigned send REST request to bittrex
func RequestSigned(method string, endpoint string, params *map[string]string, content string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	contentByte := []byte(content)

	req.Header.SetMethod(method)
	req.Header.Set("X-MBX-APIKEY", APIKey)

	req.SetRequestURI(baseURL + endpoint)
	req.SetBody(contentByte)

	for key, val := range *params {
		req.URI().QueryArgs().Add(key, val)
	}

	// START SIGNING
	raw := string(req.URI().QueryArgs().QueryString()) + string(req.Body())

	fmt.Println("raw", raw)
	mac := hmac.New(sha256.New, SecretKeyBytes)
	_, err := mac.Write([]byte(raw))

	if err != nil {
		return nil, err
	}

	req.URI().QueryArgs().Add("signature", fmt.Sprintf("%x", (mac.Sum(nil))))

	// END SIGNING
	fmt.Println("URI", string(req.URI().FullURI()))
	err = fasthttp.Do(req, resp)
	bodyBytes := resp.Body()
	return bodyBytes, err
}
