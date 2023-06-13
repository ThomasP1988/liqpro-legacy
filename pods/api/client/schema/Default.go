package schema

import (
	"github.com/graphql-go/graphql"
)

// DefaultResponseType graphql type for api key
var DefaultResponseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "DefaultResponse",
		Fields: graphql.Fields{
			"success": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)

// DefaultResponse
type DefaultResponse struct {
	Success bool `json:"success"`
}
