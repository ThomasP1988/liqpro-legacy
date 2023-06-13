package apikey

import (
	"fmt"
	common "liqpro/pods/api/client/common"
	schema "liqpro/pods/api/client/schema"
	repositories "liqpro/shared/repositories"

	"github.com/graphql-go/graphql"
)

// ListAPIKey graphql resolver
var ListAPIKey *graphql.Field = &graphql.Field{
	Type:        graphql.NewList(schema.APIKeyType),
	Description: "Get product list",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {

		userID := string(params.Context.Value(common.UserIDKey).([]byte))
		fmt.Println("ici", userID)

		apiKeys, err := repositories.GetAPIKeyRepository().ListByUser(userID)

		if err != nil {
			return nil, err
		}

		return apiKeys, nil
	},
}
