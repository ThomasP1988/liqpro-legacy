package token

import (
	"fmt"
	common "liqpro/pods/api/client/common"
	schema "liqpro/pods/api/client/schema"
	repositories "liqpro/shared/repositories"

	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

// CreateToken graphql resolver
var CreateToken *graphql.Field = &graphql.Field{
	Type:        schema.Token,
	Description: "Create new token",
	Args:        graphql.FieldConfigArgument{},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {

		userID := params.Context.Value(common.UserIDKey).([]byte)

		token, err := repositories.GetTokenRepository().Create(string(userID))

		if err != nil {
			fmt.Println("repositories.GetTokenRepository().Create", err)
			return nil, errors.New("Error creating Token")
		}

		return token, nil
	},
}
