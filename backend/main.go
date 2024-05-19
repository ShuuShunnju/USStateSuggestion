package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var stateCollection *mongo.Collection

type State struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var stateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "State",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.String},
		"name": &graphql.Field{Type: graphql.String},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"states": &graphql.Field{
			Type: graphql.NewList(stateType),
			Args: graphql.FieldConfigArgument{
				"search": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var filter bson.M
				if search, isOK := p.Args["search"].(string); isOK {
					filter = bson.M{"name": bson.M{"$regex": search, "$options": "i"}}
				} else {
					filter = bson.M{}
				}

				var states []State
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				cursor, err := stateCollection.Find(ctx, filter)
				if err != nil {
					return nil, err
				}
				defer cursor.Close(ctx)

				uniqueStates := make(map[string]bool)
				for cursor.Next(ctx) {
					var state State
					cursor.Decode(&state)
					if !uniqueStates[state.Name] {
						uniqueStates[state.Name] = true
						states = append(states, state)
					}
				}
				return states, nil
			},
		},
	},
})

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	stateCollection = client.Database("statesdb").Collection("states")

	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/graphql", gin.WrapH(handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})))
	router.Run(":8080")
}
