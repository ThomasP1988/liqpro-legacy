package apikey

import (
	"fmt"

	common "liqpro/pods/api/client/common"

	schema "liqpro/pods/api/client/schema"
	repositories "liqpro/shared/repositories"
	entities "liqpro/shared/repositories/entities"

	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

// DeleteAPIKey graphql resolver
var DeleteAPIKey *graphql.Field = &graphql.Field{
	Type:        schema.DefaultResponseType,
	Description: "Delete api key",
	Args: graphql.FieldConfigArgument{
		"accessKey": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		accessKey, _ := params.Args["accessKey"].(string)
		userID := params.Context.Value(common.UserIDKey).([]byte)

		apiKey := &entities.APIKey{}

		err := repositories.GetAPIKeyRepository().FindOne(accessKey, apiKey)

		if err != nil {
			fmt.Println("repositories.GetAPIKeyRepository().FindOne", err)
			return nil, errors.New("Error finding API Key")
		}

		if apiKey.UserID != string(userID) {
			return nil, errors.New("Unauthorised")
		}

		err = repositories.GetAPIKeyRepository().Delete(accessKey)

		if err != nil {
			fmt.Println("repositories.GetAPIKeyRepository().Delete", err)
			return nil, errors.New("Error deleting API Key")
		}

		response := schema.DefaultResponse{
			Success: true,
		}

		return response, nil
	},
}
