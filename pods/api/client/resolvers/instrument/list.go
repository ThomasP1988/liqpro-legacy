package instrument

import (
	schema "liqpro/pods/api/client/schema"

	"liqpro/config"

	"github.com/graphql-go/graphql"
)

// ListOrders graphql resolver
var List *graphql.Field = &graphql.Field{
	Type:        graphql.NewList(schema.Instrument),
	Description: "List transaction in (orders)",
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {

		var intrumentList []Instrument

		for instr, levelArr := range config.AuthorisedInstrumentsAndLevelsArray {
			intrumentList = append(intrumentList, Instrument{
				Pair:   instr,
				Levels: levelArr,
			})
		}

		return intrumentList, nil
	},
}

type Instrument struct {
	Pair   string    `json:"pair"`
	Levels []float64 `json:"levels"`
}
