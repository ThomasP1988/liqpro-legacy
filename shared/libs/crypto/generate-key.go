package crypto

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
)

// GenerateStringBytes generate random string
func GenerateStringBytes(size int) (*[]byte, error) {
	randomBytes := make([]byte, size, size)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	encodedData := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, encodedData)
	defer encoder.Close()
	encoder.Write(randomBytes)
	return &randomBytes, nil
}

// GenerateString generate random string
func GenerateString(size int) (*string, error) {
	randomBytes := make([]byte, size)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	randomString := base64.URLEncoding.EncodeToString(randomBytes)
	return &randomString, nil
}

// GenerateAccessKeyAndSecretKey generate access and secret
func GenerateAccessKeyAndSecretKey() (*string, *string, error) {
	secret, err1 := GenerateString(16)
	access, err2 := GenerateString(32)

	if err1 != nil {
		return nil, nil, err1
	}

	if err2 != nil {
		return nil, nil, err2
	}

	return secret, access, nil
}
