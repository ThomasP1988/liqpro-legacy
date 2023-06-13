package portfolio

import (
	common "liqpro/pods/api/client/common"
	schema "liqpro/pods/api/client/schema"
	repositories "liqpro/shared/repositories"

	"github.com/graphql-go/graphql"
)

// ListHold graphql resolver
var ListHold *graphql.Field = &graphql.Field{
	Type:        graphql.NewList(schema.Portfolio),
	Description: "Get portfolio",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {

		userID := string(params.Context.Value(common.UserIDKey).([]byte))
		portfolio, err := repositories.GetPortfolioRepository().ListByUser(userID)

		if err != nil {
			return nil, err
		}

		return portfolio, nil
	},
}
