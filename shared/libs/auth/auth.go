package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"liqpro/shared/repositories/cache"

	"github.com/valyala/fasthttp"

	repositories "liqpro/shared/repositories"
	entities "liqpro/shared/repositories/entities"
)

var (
	tokenNonce      []byte = []byte("Token")
	headerNonce     []byte = []byte("Nonce")
	headerSignature []byte = []byte("Signature")
	headerAPIKey    []byte = []byte("API-Key")
)

//GetUserIDFromQuery auth with url params without cache
func GetUserIDFromQuery(ctx *fasthttp.RequestCtx) (*cache.UserDataCache, error) {
	// identifying with one time token (from front end UI for example)
	token := ctx.Request.URI().QueryArgs().PeekBytes(tokenNonce)

	if len(token) > 0 {
		userID, err := repositories.GetTokenRepository().Use(string(token))

		if err != nil {
			return nil, errors.New("Wrong token.")
		}

		userData := &cache.UserDataCache{
			UserID: []byte(*userID),
		}

		return userData, nil
	}

	// identifying with api key
	nonce := ctx.Request.URI().QueryArgs().PeekBytes(headerNonce)
	signature := ctx.Request.URI().QueryArgs().PeekBytes(headerSignature)
	APIKey := ctx.Request.URI().QueryArgs().PeekBytes(headerAPIKey)

	entityAPIKey := &entities.APIKey{}

	err := repositories.GetAPIKeyRepository().FindOne(string(APIKey), entityAPIKey)

	if err != nil {
		return nil, errors.New("API Key not found")
	}

	nonceLength := len(nonce)

	if nonceLength < 32 || nonceLength > 65 {
		return nil, errors.New("Your Nonce should have at least 32 characters and 64 characters maximum")
	}

	if !ValidMAC(nonce, signature, []byte(entityAPIKey.SecretKey)) {
		return nil, errors.New("Your signature is incorrect")
	}

	return cache.CreateUserDataCache(APIKey)
}

// GetUserDataWithHeadersFromCache verify signature and return data from cache
func GetUserDataWithHeadersFromCache(ctx *fasthttp.RequestCtx) (*cache.UserDataCache, error) {
	var (
		userData *cache.UserDataCache
		err      error
		apiKey   []byte = ctx.Request.Header.PeekBytes(headerAPIKey)
	)

	userData, err = cache.GetUserByAPIKey(apiKey)

	if err != nil {
		userData, err = cache.SetUserByAPIKey(apiKey)
		if err != nil {
			return nil, err
		}
	}

	nonceLength := len(ctx.Request.Header.PeekBytes(headerNonce))

	if nonceLength < 32 || nonceLength > 65 {
		return nil, errors.New("Your Nonce should have at least 32 characters and 64 characters maximum")
	}

	if !ValidMAC(ctx.Request.Header.PeekBytes(headerNonce), ctx.Request.Header.PeekBytes(headerSignature), userData.SecretKey) {
		return nil, errors.New("Your signature is incorrect")
	}

	return userData, nil
}

// ValidMAC Package hmac implements the Keyed-Hash Message Authentication Code (HMAC) as defined in U.S. Federal Information Processing Standards Publication 198.
func ValidMAC(message, messageMAC, key []byte) bool {
	var expectedMAC []byte = make([]byte, len(messageMAC))
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	base64.StdEncoding.Encode(expectedMAC, mac.Sum(nil))
	return hmac.Equal(messageMAC, expectedMAC)
}
