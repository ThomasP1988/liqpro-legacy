package schema

import (
	"github.com/graphql-go/graphql"
)

// Token graphql type for api key
var Token = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Token",
		Fields: graphql.Fields{
			"value": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
