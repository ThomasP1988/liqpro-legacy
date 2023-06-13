package apikey

import (
	"fmt"
	common "liqpro/pods/api/client/common"
	schema "liqpro/pods/api/client/schema"
	crypto "liqpro/shared/libs/crypto"
	repositories "liqpro/shared/repositories"
	entities "liqpro/shared/repositories/entities"

	"time"

	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

// CreateAPIKey graphql resolver
var CreateAPIKey *graphql.Field = &graphql.Field{
	Type:        schema.APIKeyType,
	Description: "Create new api key",
	Args:        graphql.FieldConfigArgument{},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		accessKey, secretKey, err := crypto.GenerateAccessKeyAndSecretKey()

		if err != nil {
			return nil, err
		}

		userID := params.Context.Value(common.UserIDKey).([]byte)

		apiKey := &entities.APIKey{
			AccessKey: *accessKey,
			SecretKey: *secretKey,
			UserID:    string(userID),
			Created:   time.Now().Unix(),
		}
		err = repositories.GetAPIKeyRepository().Create(apiKey)

		if err != nil {
			fmt.Println("repositories.GetAPIKeyRepository().Create", err)
			return nil, errors.New("Error creating API Key")
		}

		return apiKey, nil
	},
}
