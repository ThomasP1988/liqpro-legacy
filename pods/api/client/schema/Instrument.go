package schema

import (
	"github.com/graphql-go/graphql"
)

// Order graphql type for api key
var Instrument = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Instrument",
		Fields: graphql.Fields{
			"pair": &graphql.Field{
				Type: graphql.String,
			},
			"levels": &graphql.Field{
				Type: graphql.NewList(graphql.Float),
			},
		},
	},
)
