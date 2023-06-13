package schema

import (
	"github.com/graphql-go/graphql"
)

// Portfolio graphql type for api key
var Portfolio = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Hold",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"currency": &graphql.Field{
				Type: graphql.String,
			},
			"userId": &graphql.Field{
				Type: graphql.String,
			},
			"total": &graphql.Field{
				Type: graphql.Float,
			},
			"lastModified": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)
