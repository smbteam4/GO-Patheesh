// this package can be included api handler related functionalities
package handlers

import (
	"fmt"
	"robot-apocalypse/pkg/db"
	"robot-apocalypse/pkg/models"

	"go.uber.org/zap"
)

// handler struct
type Handler struct {
	Logger *zap.Logger
	DB     *db.MongoAdapter
}

// initiate new handler
func NewHandler(logger *zap.Logger, db *db.MongoAdapter) *Handler {
	return &Handler{
		Logger: logger,
		DB:     db,
	}
}

// new survivor handler
// create new survivor entry to the database
func (handle *Handler) NewSurvivorHandler(sr models.Survivor) error {
	//check user id already exists
	exists, err := handle.DB.Survivors().CheckSurvivorExists(sr.ID)
	if err != nil {
		return fmt.Errorf("unable to process your request")
	}
	if exists {
		return fmt.Errorf("survivor entry already exists")
	}

	// create new survivor
	return handle.DB.Survivors().New(sr)
}

// new survivor handler
// create new survivor entry to the database
func (handle *Handler) UpdateSurvivorHandler(sr models.Survivor) error {
	//check user id already exists
	exists, err := handle.DB.Survivors().CheckSurvivorExists(sr.ID)
	if err != nil {
		return fmt.Errorf("unable to process your request")
	}
	if !exists {
		return fmt.Errorf("survivor not exists in the system")
	}

	// update the survivor
	return handle.DB.Survivors().Update(sr)
}

// mark a survivor as infected
func (handle *Handler) MarkSurvivorInfectedHandler(sr models.SurvivorInfected) error {
	//check user id already exists
	exists, err := handle.DB.Survivors().CheckSurvivorExists(sr.ID)
	if err != nil {
		return fmt.Errorf("unable to process your request")
	}
	if !exists {
		return fmt.Errorf("unable to identify the survivor")
	}

	// mark the survivor as infetcted
	return handle.DB.Survivors().Infected(sr.ID, sr.ReportedBy)
}

// infected/ non infected percentage
func (handle *Handler) InfectionPercentagehandler() (*models.InfectionReport, error) {
	infectedCount, err := handle.DB.Survivors().InfectedCount()
	if err != nil {
		return nil, fmt.Errorf("unable to process your request")
	}
	totalSurvivors, err := handle.DB.Survivors().TotalSurvivors()
	if err != nil {
		return nil, fmt.Errorf("unable to process your request")
	}

	return &models.InfectionReport{
		Infected:    (float32(infectedCount) * 100) / float32(totalSurvivors),
		NonInfected: (float32(totalSurvivors-infectedCount) * 100) / float32(totalSurvivors),
	}, nil
}

// list our infected or non infected survivors list
func (handle *Handler) InfectionNonInfectionListhandler(criteria string) ([]models.Survivor, error) {
	if criteria != "infected" && criteria != "non-infected" {
		return nil, fmt.Errorf("invalid url %v", criteria)
	}
	data, err := handle.DB.Survivors().GetSurvivors(criteria)
	if err != nil {
		return nil, fmt.Errorf("unable to process your request")
	}

	return data, nil
}

// list our infected or non infected survivors list
func (handle *Handler) LoadRobotsHandler(robotList []models.RobotList) error {
	err := handle.DB.Robots().LoadData(robotList)
	if err != nil {
		return fmt.Errorf("unable to process your request")
	}
	return nil
}

// list our infected or non infected survivors list
func (handle *Handler) ListRobotsHandler() ([]models.RobotList, error) {
	robotList, err := handle.DB.Robots().ListData()
	if err != nil {
		return nil, fmt.Errorf("unable to process your request")
	}
	return robotList, nil
}
