package schema

import (
	"github.com/graphql-go/graphql"
)

// APIKeyType graphql type for api key
var APIKeyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "APIKey",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"accessKey": &graphql.Field{
				Type: graphql.String,
			},
			"userId": &graphql.Field{
				Type: graphql.String,
			},
			"secretKey": &graphql.Field{
				Type: graphql.String,
			},
			"created": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)
