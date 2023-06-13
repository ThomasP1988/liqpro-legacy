package cache

import (
	"errors"
	"fmt"
	repositories "liqpro/shared/repositories"
	entities "liqpro/shared/repositories/entities"

	"github.com/dgraph-io/ristretto"
)

// APIUsers to store user data
var APIUsers *ristretto.Cache

// UserDataCache data we store into the cache
type UserDataCache struct {
	SecretKey []byte
	UserID    []byte
}

// InitAPIUsers init the cache
func InitAPIUsers() error {
	// cache
	var errCache error
	APIUsers, errCache = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if errCache != nil {
		return errCache
	}
	return nil
}

// GetUserByAPIKey retrieve user data with APIKey
func GetUserByAPIKey(APIKey []byte) (*UserDataCache, error) {
	value, exist := APIUsers.Get(APIKey)
	fmt.Println("value", value)
	if !exist {
		return nil, errors.New("User not in cache")
	}

	valueTyped := value.(UserDataCache)

	return &valueTyped, nil
}

// CreateUserDataCache create data structure to store the data we need about the user
func CreateUserDataCache(APIKey []byte) (*UserDataCache, error) {
	userDataCache := &UserDataCache{}
	entityAPIKey := &entities.APIKey{}

	err := repositories.GetAPIKeyRepository().FindOne(string(APIKey), entityAPIKey)

	if err != nil {
		return nil, errors.New("API Key not found")
	}

	userDataCache.UserID = []byte(entityAPIKey.UserID)
	userDataCache.SecretKey = []byte(entityAPIKey.SecretKey)
	return userDataCache, nil
}

// SetUserByAPIKey retrieve user data with APIKey
func SetUserByAPIKey(APIKey []byte) (*UserDataCache, error) {

	userDataCache, err := CreateUserDataCache(APIKey)

	if err != nil {
		return nil, err
	}

	APIUsers.Set(APIKey, *userDataCache, 1)

	return userDataCache, nil
}
