// package database
// This package will be included all the data, methods and implementations related to
// the database
package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// mongo connection adapter
type MongoAdapter struct {
	client   *mongo.Client
	database *mongo.Database

	Services
}

// adapter services
type Services interface {
	Survivors() SurvivorServices
	Robots() RobotsServices
}

// survivor service
func (adptr *MongoAdapter) Survivors() *SurvivorServices {
	srv := NewSurvivorServices()
	srv.Collection = adptr.ConnectCollection("survivors")
	srv.LocationHistory = adptr.ConnectCollection("survivors_location_history")
	return srv
}

// survivor service
func (adptr *MongoAdapter) Robots() *RobotsServices {
	srv := NewRobotsServices()
	srv.Collection = adptr.ConnectCollection("robots")
	return srv
}

// Connect to cllection
// Create a handle to the respective collection in the database.
func (mongoadapter *MongoAdapter) ConnectCollection(tb string) *mongo.Collection {
	collection := mongoadapter.database.Collection(tb)
	return collection
}

// initiate new mongo db connection
func NewConnection(ctx context.Context, URI string, database string) (*MongoAdapter, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		return nil, err
	}

	// connect to the database
	db := client.Database(database)

	// and initiate a ping request to check database is alive
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("could not ping to mongo db service: %v", err)
	}

	adptr := &MongoAdapter{
		client:   client,
		database: db,
	}

	return adptr, nil
}
