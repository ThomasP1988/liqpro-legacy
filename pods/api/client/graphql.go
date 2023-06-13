package main

import (
	"context"
	"fmt"

	common "liqpro/pods/api/client/common"
	apikey "liqpro/pods/api/client/resolvers/apikey"
	"liqpro/pods/api/client/resolvers/instrument"
	"liqpro/pods/api/client/resolvers/order"
	"liqpro/pods/api/client/resolvers/portfolio"
	"liqpro/pods/api/client/resolvers/token"

	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"listApiKeys":    apikey.ListAPIKey,
			"portfolio":      portfolio.ListHold,
			"orders":         order.ListOrders,
			"getToken":       token.CreateToken,
			"listIntruments": instrument.List,
			/* Get (read) single product by id
			   http://localhost:8080/product?query={product(id:1){name,info,price}}
			*/
			// "product": &graphql.Field{
			// 	Type:        productType,
			// 	Description: "Get product by id",
			// 	Args: graphql.FieldConfigArgument{
			// 		"id": &graphql.ArgumentConfig{
			// 			Type: graphql.Int,
			// 		},
			// 	},
			// 	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 		id, ok := p.Args["id"].(int)
			// 		if ok {
			// 			// Find product
			// 			for _, product := range products {
			// 				if int(product.ID) == id {
			// 					return product, nil
			// 				}
			// 			}
			// 		}
			// 		return nil, nil
			// 	},
			// },
			/* Get (read) product list
			   http://localhost:8080/product?query={list{id,name,info,price}}
			*/
			// "list": &graphql.Field{
			// 	Type:        graphql.NewList(productType),
			// 	Description: "Get product list",
			// 	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// 		return products, nil
			// 	},
			// },
		},
	})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createApiKey": apikey.CreateAPIKey,
		"deleteApiKey": apikey.DeleteAPIKey,
		/* Create new product item
		http://localhost:8080/product?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}
		*/
		// "create": &graphql.Field{
		// 	Type:        productType,
		// 	Description: "Create new product",
		// 	Args: graphql.FieldConfigArgument{
		// 		"name": &graphql.ArgumentConfig{
		// 			Type: graphql.NewNonNull(graphql.String),
		// 		},
		// 		"info": &graphql.ArgumentConfig{
		// 			Type: graphql.String,
		// 		},
		// 		"price": &graphql.ArgumentConfig{
		// 			Type: graphql.NewNonNull(graphql.Float),
		// 		},
		// 	},
		// 	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		// 		rand.Seed(time.Now().UnixNano())
		// 		product := Product{
		// 			ID:    int64(rand.Intn(100000)), // generate random ID
		// 			Name:  params.Args["name"].(string),
		// 			Info:  params.Args["info"].(string),
		// 			Price: params.Args["price"].(float64),
		// 		}
		// 		products = append(products, product)
		// 		return product, nil
		// 	},
		// },
	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

func executeQuery(p *PostData, c *fiber.Ctx) *graphql.Result {

	ctx := context.WithValue(context.Background(), common.UserIDKey, c.Request().Header.Peek("userId"))

	result := graphql.Do(graphql.Params{
		Context:        ctx,
		Schema:         schema,
		RequestString:  p.Query,
		VariableValues: p.Variables,
		OperationName:  p.Operation,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

// PostData data we get from GraphQL Front end
type PostData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}
