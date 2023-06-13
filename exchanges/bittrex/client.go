package bittrex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	apiKey    = ""
	apiSecret = ""
	baseURL   = "https://api.bittrex.com/v3"
)

// Request send REST request to bittrex
func Request(method string, endpoint string, params map[string]string, content *[]byte) (*[]byte, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	for key, val := range params {
		req.URI().QueryArgs().Add(key, val)
	}
	req.Header.SetContentType("application/json")
	// req.Header.Add("content-type", "application/json")
	req.Header.Add("Api-Key", apiKey)

	// timestamp
	t := strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
	req.Header.Add("Api-Timestamp", t)

	if content == nil {
		content = &[]byte{}
	}

	// content SHA512 to hex encoding
	contentSum := sha512.Sum512(*content)
	contentHash := hex.EncodeToString(contentSum[:])

	req.Header.Add("Api-Content-Hash", contentHash)

	// req.Header.Add("Api-Subaccount-Id", API_KEY)

	// signature
	req.Header.Add("Api-Signature", encodeSignature(t, baseURL+endpoint, method, contentHash))
	fmt.Println(encodeSignature(t, baseURL+endpoint, method, contentHash))

	req.Header.SetMethod(method)
	req.SetRequestURI(baseURL + endpoint)
	req.SetBody(*content)

	err := fasthttp.Do(req, resp)
	bodyBytes := resp.Body()
	fmt.Println(resp.StatusCode())
	resp.StatusCode()
	return &bodyBytes, err
}

func encodeSignature(t string, URI string, method string, contentHash string) string {
	toSign := []string{t, URI, method, contentHash}
	sig := hmac.New(sha512.New, []byte(apiSecret))
	sig.Write([]byte(strings.Join(toSign, "")))
	return hex.EncodeToString(sig.Sum(nil))
}
