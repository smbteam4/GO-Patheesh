package db

import (
	"context"
	"robot-apocalypse/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var InfectionMinimumReportCount = 3

// survivor services
type SurvivorServices struct {
	Collection      *mongo.Collection
	LocationHistory *mongo.Collection
}

// initiate new survivor services
func NewSurvivorServices() *SurvivorServices {
	return &SurvivorServices{}
}

// New survivor entry
func (sr *SurvivorServices) New(data models.Survivor) error {
	_, err := sr.Collection.InsertOne(context.TODO(), data)
	return err
}

// fetch the survivor details
func (sr *SurvivorServices) GetSurvivor(id string) (*models.Survivor, error) {
	var collected_data *models.Survivor
	queryOptions := options.FindOneOptions{}

	err := sr.Collection.FindOne(context.TODO(), bson.M{
		"id": id,
	}, &queryOptions).Decode(&collected_data)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	return collected_data, nil
}

// check survivor exists or not with the id
func (sr *SurvivorServices) CheckSurvivorExists(id string) (bool, error) {
	exists, err := sr.Collection.CountDocuments(context.TODO(), bson.M{
		"id": id,
	})

	if err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}

	// if the entry exists, return true
	if exists > 0 {
		return true, nil
	}
	// if entry not exists
	return false, nil
}

// Update survivor entry
func (sr *SurvivorServices) Update(data models.Survivor) error {
	var currentLocation models.Location

	// based on the input it will update the entries
	updateList := bson.M{}
	if data.Name != "" {
		updateList["name"] = data.Name
	}
	if data.Age > 0 {
		updateList["age"] = data.Age
	}

	if (data.Location != models.Location{}) {
		updateList["location"] = data.Location
		currentLocation = data.Location
	}

	// update
	_, err := sr.Collection.UpdateOne(
		context.TODO(),
		bson.M{
			"id": data.ID,
		},
		bson.M{
			"$set": updateList,
		})

	// when the request includes location change, then we will keep an history.
	// with that history we can locate the escape track of that survivor
	if (currentLocation != models.Location{}) {
		sr.NewLocationHistory(data.ID, data.Location)
	}

	return err
}

// Update survivor entry
func (sr *SurvivorServices) Infected(id string, infect_reported string) error {
	_, err := sr.Collection.UpdateOne(context.TODO(), bson.M{
		"id": id,
	}, bson.M{
		"$inc": bson.M{
			"reportedcount": 1,
		},
		"$push": bson.M{
			"reportedby": infect_reported,
		},
	})
	return err
}

// insert new change location history
func (sr *SurvivorServices) NewLocationHistory(id string, location models.Location) error {
	_, err := sr.LocationHistory.InsertOne(context.TODO(), bson.M{
		"id":       id,
		"location": location,
	})

	return err
}

// prepare infection report
func (sr *SurvivorServices) InfectedCount() (int, error) {
	count, err := sr.Collection.CountDocuments(context.TODO(), bson.M{
		"reportedcount": bson.M{"$gte": InfectionMinimumReportCount},
	})

	if err != nil && err != mongo.ErrNoDocuments {
		return 0, err
	}

	return int(count), nil
}

// Count total survivors
func (sr *SurvivorServices) TotalSurvivors() (int, error) {
	count, err := sr.Collection.CountDocuments(context.TODO(), bson.M{})

	if err != nil && err != mongo.ErrNoDocuments {
		return 0, err
	}

	return int(count), nil
}

// list infected survivors
func (sr *SurvivorServices) GetSurvivors(criteria string) ([]models.Survivor, error) {
	var collected_data []models.Survivor
	ctx := context.TODO()
	var filter bson.M
	if criteria == "infected" {
		filter = bson.M{
			"reportedcount": bson.M{"$gte": InfectionMinimumReportCount},
		}
	} else {
		filter = bson.M{
			"reportedcount": bson.M{"$lt": InfectionMinimumReportCount},
		}
	}

	cursor, err := sr.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		// Declare a result BSON object
		var result models.Survivor
		err = cursor.Decode(&result)
		if err != nil {
			continue
		}
		collected_data = append(collected_data, result)
	}
	return collected_data, nil
}
