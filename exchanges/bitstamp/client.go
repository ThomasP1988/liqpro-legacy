package bitstamp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	baseURL = "https://www.bitstamp.net/api/v2"
)

// ClientID Bitstamp API credentials
const ClientID string = ""

// APIKey Bitstamp API credentials
const APIKey string = ""

// SecretKey Bitstamp API credentials
const SecretKey string = ""

// ClientIDBytes Bitstamp API credentials
var ClientIDBytes = []byte(ClientID)

// APIKeyBytes Bitstamp API credentials
var APIKeyBytes = []byte(APIKey)

// SecretKeyBytes Bitstamp API credentials
var SecretKeyBytes = []byte(SecretKey)

var urlFormEncoded = []byte("application/x-www-form-urlencoded")

// Request send REST request to bitstamp
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
	resp.StatusCode()
	return bodyBytes, err
}

// RequestSigned send REST request to bittrex
func RequestSigned(method string, endpoint string, params url.Values) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	// add required key, signature & nonce to values

	nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
	mac := hmac.New(sha256.New, SecretKeyBytes)
	fmt.Println("nonce", nonce)

	mac.Write([]byte(nonce + ClientID + APIKey))

	params.Set("key", APIKey)
	params.Set("signature", strings.ToUpper(hex.EncodeToString(mac.Sum(nil))))
	params.Set("nonce", nonce)

	req.Header.SetMethod(method)
	req.Header.SetContentTypeBytes(urlFormEncoded)
	req.SetRequestURI(baseURL + endpoint)
	req.SetBody([]byte(params.Encode()))

	err := fasthttp.Do(req, resp)
	bodyBytes := resp.Body()
	resp.StatusCode()
	return bodyBytes, err
}
