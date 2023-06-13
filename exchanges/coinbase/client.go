package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

// PassPhrase Coinbase API credentials
const PassPhrase string = "Ici c'est Paris."

// APIKeyBytes Coinbase API credentials
var APIKeyBytes = []byte(APIKey)

// SecretKeyBytes Coinbase API credentials
var SecretKeyBytes = []byte(SecretKey)

// PassPhraseBytes Coinbase API credentials
var PassPhraseBytes = []byte(PassPhrase)

const baseURL = "https://api-public.pro.coinbase.com"
const baseURLTest = "https://api-public.sandbox.pro.coinbase.com"

var headerApplicationJSON []byte = []byte("application/json")
var headerAccept []byte = []byte("Accept")
var headerAccessKey []byte = []byte("CB-ACCESS-KEY")
var headerAccessPassPhrase []byte = []byte("CB-ACCESS-PASSPHRASE")
var headerAccessTimestamp []byte = []byte("CB-ACCESS-TIMESTAMP")
var headerSign []byte = []byte("CB-ACCESS-SIGN")

// RequestSigned send REST request to bittrex
func RequestSigned(method string, endpoint string, body *[]byte) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	// add required key, signature & nonce to values
	req.Header.SetContentTypeBytes(headerApplicationJSON)
	req.Header.AddBytesKV(headerAccept, headerApplicationJSON)

	req.Header.SetMethod(method)
	req.SetRequestURI(baseURLTest + endpoint)
	AddHeaders(req, method, endpoint, string(*body))
	req.SetBody(*body)

	err := fasthttp.Do(req, resp)
	bodyBytes := resp.Body()
	return bodyBytes, err
}

// AddHeaders to signed request
func AddHeaders(req *fasthttp.Request, method, url, data string) error {

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	req.Header.AddBytesKV(headerAccessKey, APIKeyBytes)
	req.Header.AddBytesKV(headerAccessPassPhrase, PassPhraseBytes)
	req.Header.AddBytesK(headerAccessTimestamp, timestamp)

	message := timestamp + method + url + data

	sig, err := generateSig(message, SecretKey)
	if err != nil {
		return err
	}
	req.Header.AddBytesK(headerSign, sig)

	return nil
}

func generateSig(message, secret string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	signature := hmac.New(sha256.New, key)
	_, err = signature.Write([]byte(message))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature.Sum(nil)), nil
}
