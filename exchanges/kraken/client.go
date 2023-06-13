package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

// APIKey Kraken API credentials
const APIKey string = ""

// SecretKey Kraken API credentials
const SecretKey string = ""

// APIKeyBytes Kraken API credentials
var APIKeyBytes = []byte(APIKey)

// SecretKeyBytes Kraken API credentials
var SecretKeyBytes = []byte(SecretKey)

// change to https://api-aws.kraken.pro if hosted in AWS Japan https://krakenapi.github.io/docs/spot/v1/en/#access-urls
const (
	baseURL = "https://api.kraken.com"
	host    = "api.kraken.com"
	// APIVersion - Kraken API Version Number
	APIVersion = "0"
)

var applicationJSON []byte = []byte("application/x-www-form-urlencoded; charset=utf-8")
var headerKeyAPIKey []byte = []byte("API-Key")
var headerKeyAPISign []byte = []byte("API-Sign")
var headerUserAgent []byte = []byte("User-Agent")
var methodPOST []byte = []byte("POST")

// RequestSigned send REST request to bittrex
func RequestSigned(endpoint string, params *url.Values) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	req.Header.SetContentTypeBytes(applicationJSON)
	req.Header.SetMethodBytes(methodPOST)

	requestURL := baseURL + "/" + APIVersion + "/private/" + endpoint

	req.SetRequestURI(requestURL)
	addHeaders(req, params, requestURL)
	req.AppendBodyString(params.Encode())
	println("req", req.String())
	err := fasthttp.Do(req, resp)

	return resp.Body(), err
}

func addHeaders(req *fasthttp.Request, params *url.Values, url string) error {
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	params.Set("nonce", nonce)
	strings.NewReader(params.Encode())

	req.Header.AddBytesKV(headerKeyAPIKey, APIKeyBytes)

	signature, err := sign(params, url)
	if err != nil {
		fmt.Println("error signing kraken request", err)
		return err
	}
	req.Header.AddBytesK(headerKeyAPISign, signature)
	return nil
}

func sign(params *url.Values, url string) (string, error) {
	sha := sha256.New()

	if _, err := sha.Write([]byte(params.Get("nonce") + params.Encode())); err != nil {
		return "", err
	}
	hashData := sha.Sum(nil)
	s, err := base64.StdEncoding.DecodeString(SecretKey)
	if err != nil {
		return "", err
	}
	hmacObj := hmac.New(sha512.New, s)

	if _, err := hmacObj.Write(append([]byte(url), hashData...)); err != nil {
		return "", err
	}
	hmacData := hmacObj.Sum(nil)
	return base64.StdEncoding.EncodeToString(hmacData), nil
}
