package huobi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

// AccountID Huobi Credential
const AccountID string = "36570531"

// APIKey Huobi API credentials
const APIKey string = ""

// SecretKey Huobi API credentials
const SecretKey string = ""

// APIKeyBytes Huobi API credentials
var APIKeyBytes = []byte(APIKey)

// SecretKeyBytes Huobi API credentials
var SecretKeyBytes = []byte(SecretKey)

// change to https://api-aws.huobi.pro if hosted in AWS Japan https://huobiapi.github.io/docs/spot/v1/en/#access-urls
const (
	baseURL = "https://api.huobi.pro"
	host    = "api.huobi.pro"
)

var applicationJSON []byte = []byte("application/json")

// RequestSigned send REST request to bittrex
func RequestSigned(method string, endpoint string, content []byte, params *url.Values) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	req.Header.SetContentTypeBytes(applicationJSON)
	req.Header.SetMethod(method)

	if params == nil {
		params = &url.Values{}
	}

	fmt.Println("url", GetURL(method, endpoint, params))
	// fmt.Println("content", string(content))

	req.SetRequestURI(GetURL(method, endpoint, params))
	if content != nil {
		req.SetBody(content)
	}

	err := fasthttp.Do(req, resp)

	return resp.Body(), err
}

// GetURL to signed request
func GetURL(method string, endpoint string, params *url.Values) string {

	now := time.Now().UTC()
	time := now.Format("2006-01-02T15:04:05")
	params.Set("AccessKeyId", APIKey)
	params.Set("SignatureMethod", "HmacSHA256")
	params.Set("SignatureVersion", "2")
	params.Set("Timestamp", time)
	encodedParams := params.Encode()
	signature := generateSig(method, endpoint, encodedParams)

	url := fmt.Sprintf("https://%s%s?%s&Signature=%s", host, endpoint, encodedParams, url.QueryEscape(signature))

	return url

}

func generateSig(method string, path string, parameters string) string {

	var sb strings.Builder
	sb.WriteString(method)
	sb.WriteString("\n")
	sb.WriteString(host)
	sb.WriteString("\n")
	sb.WriteString(path)
	sb.WriteString("\n")
	sb.WriteString(parameters)

	return sign(sb.String())
}

func sign(payload string) string {
	hash := hmac.New(sha256.New, SecretKeyBytes)
	hash.Write([]byte(payload))
	result := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return result
}
