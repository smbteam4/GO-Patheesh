// package model
// This package will include the model structure
package models

// env configruation
type EnvironmentalConfigs struct {
	ServerPort    string `default:"8080" split_words:"true"`
	Mode          string `default:"debug" split_words:"true"`
	MongoHost     string `default:"mongodb://admin:admin@localhost:27017/" split_words:"true"`
	MongoDatabase string `default:"robot-apocalypse" split_words:"true"`
}

// model survivor
type Survivor struct {
	// id
	ID string `json:"id"`
	// name
	Name string `json:"name"`
	// age
	Age int `json:"age"`
	// location
	Location Location `json:"location"`
	// resources
	Resources Resources `json:"resources"`
	// reported count
	ReportedCount int `json:"reportedcount,omitempty"`
}

// model location
// Which can be hold the location related information
type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// model Resources
// inventory of resources
type Resources []string

// model response
// global client response structure
// swagger:model APIResponseModel
type APIResponse struct {
	StatusCode int         `json:"status_code"`       // status code
	Message    string      `json:"message,omitempty"` // response message
	Data       interface{} `json:"data,omitempty"`    // response data
}

// model survivor infected
type SurvivorInfected struct {
	// id
	ID string `json:"id"`
	// reported_by
	ReportedBy string `json:"reported_by"`
}

// infection report
type InfectionReport struct {
	// infected
	Infected float32 `json:"infected"`
	// non_infected
	NonInfected float32 `json:"non_infected"`
}

// robot list
type RobotList struct {
	Model            string `json:"model"`
	SerialNumber     string `json:"serialNumber"`
	ManufacturedDate string `json:"manufacturedDate"`
	Category         string `json:"category"`
}
