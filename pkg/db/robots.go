package db

import (
	"context"
	"robot-apocalypse/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// robots services
type RobotsServices struct {
	Collection *mongo.Collection
}

// initiate new service
func NewRobotsServices() *RobotsServices {
	return &RobotsServices{}
}

// New survivor entry
func (sr *RobotsServices) LoadData(data []models.RobotList) error {
	docs := make([]interface{}, 0)
	for _, doc := range data {
		docs = append(docs, doc)
	}
	// remove existing lists and store new items
	sr.Collection.DeleteMany(context.TODO(), bson.M{})

	_, err := sr.Collection.InsertMany(context.TODO(), docs)
	return err
}

// list robots
func (sr *RobotsServices) ListData() ([]models.RobotList, error) {
	var collected_data []models.RobotList
	ctx := context.TODO()
	cursor, err := sr.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		// Declare a result BSON object
		var result models.RobotList
		err = cursor.Decode(&result)
		if err != nil {
			continue
		}
		collected_data = append(collected_data, result)
	}
	return collected_data, nil
}
