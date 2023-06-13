package order

import (
	common "liqpro/pods/api/client/common"
	schema "liqpro/pods/api/client/schema"
	repositories "liqpro/shared/repositories"

	"github.com/graphql-go/graphql"
)

// ListOrders graphql resolver
var ListOrders *graphql.Field = &graphql.Field{
	Type:        graphql.NewList(schema.Order),
	Description: "List transaction in (orders)",
	Args: graphql.FieldConfigArgument{
		"skip": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"limit": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {

		skip, _ := params.Args["skip"].(int)
		limit, _ := params.Args["limit"].(int)

		skip64 := int64(skip)
		limit64 := int64(limit)

		userID := string(params.Context.Value(common.UserIDKey).([]byte))
		orders, err := repositories.GetTransactionInRepository().ListByUser(userID, &skip64, &limit64)

		if err != nil {
			return nil, err
		}

		return orders, nil
	},
}
