package schema

import (
	"github.com/graphql-go/graphql"
)

// Order graphql type for api key
var Order = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Order",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"action": &graphql.Field{
				Type: graphql.Int,
			},
			"instrument": &graphql.Field{
				Type: graphql.String,
			},
			"quantityAsked": &graphql.Field{
				Type: graphql.Float,
			},
			"quantityExecuted": &graphql.Field{
				Type: graphql.Float,
			},
			"pricePerUnit": &graphql.Field{
				Type: graphql.Float,
			},
			"totalPrice": &graphql.Field{
				Type: graphql.Float,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"dateExecuted": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)
