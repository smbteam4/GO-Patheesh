// Package classification  Golang Robot-Apocalypse API.
//
// Golang Robot-Apocalypse API.
//
//     Schemes: http
//     BasePath: /
//     Version: 1.0.0
//     Host: localhost:8080/api/v1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// security:
//   - APIKeyHeader: []
//
// securityDefinitions:
//  APIKeyHeader:
//    type: apiKey
//    in: header
//    name: TOKEN
//
// swagger:meta
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"robot-apocalypse/handlers"
	"robot-apocalypse/pkg/db"
	"robot-apocalypse/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

var (
	appName      = "robotApocalypse"         // application name
	logger       *zap.Logger                 // zap logger
	apiConfig    models.EnvironmentalConfigs // api environmental config
	mongoAdapter *db.MongoAdapter            // mongo connection holder
)

func main() {
	// load the environmental configurations
	err := envconfig.Process(appName, &apiConfig)
	if err != nil {
		panic(fmt.Errorf("%s: %s", appName, err))
	}
	// initiate logger
	{
		if apiConfig.Mode == "debug" {
			logger, err = zap.NewDevelopment()
		} else {
			logger, err = zap.NewProduction()
		}
		if err != nil {
			log.Fatal(err)
		}
		defer logger.Sync()
	}

	// initiate mongo db connection
	mongoAdapter, err = db.NewConnection(context.TODO(), apiConfig.MongoHost, apiConfig.MongoDatabase)
	if err != nil {
		logger.Error("unable to initialize database connection", zap.Error(err))
		return
	}

	// initiate fiber router
	app := fiber.New(fiber.Config{
		AppName: appName,
	})

	// initiate routers and handlers
	InitRouterhandlers(app)

	// listener
	logger.Info("...starting the server...")
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", apiConfig.ServerPort)); err != nil {
			logger.Error("unable to start the server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)   // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt) // When an interrupt is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown() // shutdown
}

// InitRouterhandlers
// This method id used to initiate all the api endpoints and its handler
// methods
func InitRouterhandlers(app *fiber.App) {
	// initiate new api handler object
	handler := handlers.NewHandler(logger, mongoAdapter)

	// initiate a /api/v1 endpoint
	v1 := app.Group("/api").Group("/v1")

	// swagger:route POST /survivors Survivors idOfSurvivorCreateEndpoint
	// create new survivor endpoint
	//
	// responses:
	//   200: APIResponseModel
	v1.Post("/survivors", func(c *fiber.Ctx) error {
		var survivor models.Survivor

		// parse the request body
		if err := c.BodyParser(&survivor); err != nil {
			logger.Error("unable to parse the request", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "unable to parse the request",
			})
		}

		err := handler.NewSurvivorHandler(survivor)
		if err != nil {
			logger.Error("unable to add new survivor", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		return c.Status(http.StatusOK).JSON(models.APIResponse{
			StatusCode: http.StatusOK,
			Message:    "successfully added survivor",
		})
	})

	// update endpoint
	// swagger:route PUT /survivors Survivors idOfSurvivorUpdateEndpoint
	// update the survivor informations
	//
	// responses:
	//   200: APIResponseModel
	v1.Put("/survivors", func(c *fiber.Ctx) error {
		var survivor models.Survivor

		// parse the request body
		if err := c.BodyParser(&survivor); err != nil {
			logger.Error("unable to parse the request", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "unable to parse the request",
			})
		}

		err := handler.UpdateSurvivorHandler(survivor)
		if err != nil {
			logger.Error("unable to update survivor", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		return c.Status(http.StatusOK).JSON(models.APIResponse{
			StatusCode: http.StatusOK,
			Message:    "successfully updated survivor",
		})
	})

	// mark the survivor as infected
	// swagger:route Put /survivors/infected Survivors idOfSurvivorInfectedEndpoint
	// mark the survivor as infected
	//
	// responses:
	//   200: APIResponseModel
	v1.Put("/survivors/infected", func(c *fiber.Ctx) error {
		var survivor models.SurvivorInfected

		// parse the request body
		if err := c.BodyParser(&survivor); err != nil {
			logger.Error("unable to parse the request", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "unable to parse the request",
			})
		}

		err := handler.MarkSurvivorInfectedHandler(survivor)
		if err != nil {
			logger.Error("unable to update survivor", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		return c.Status(http.StatusOK).JSON(models.APIResponse{
			StatusCode: http.StatusOK,
			Message:    "successfully reported survivor as infected",
		})
	})

	// infected percentage
	// swagger:route GET /report/percentage Report idOfreportPercentage
	// Percentage
	//
	// responses:
	//   200: APIResponseModel
	v1.Get("/report/percentage", func(c *fiber.Ctx) error {
		reportData, err := handler.InfectionPercentagehandler()
		if err != nil {
			logger.Error("unable to fetch the survivor infection report", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		return c.Status(http.StatusOK).JSON(models.APIResponse{
			StatusCode: http.StatusOK,
			Data:       reportData,
		})
	})

	// list of criteria
	// swagger:route GET /report/{criteria} Report idOfReportCriteriaEndpoint
	// Percentage
	//
	// responses:
	//   200:
	v1.Get("/report/:criteria", func(c *fiber.Ctx) error {
		reportData, err := handler.InfectionNonInfectionListhandler(c.Params("criteria"))
		if err != nil {
			logger.Error("unable to fetch the survivor infection list", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		return c.Status(http.StatusOK).JSON(models.APIResponse{
			StatusCode: http.StatusOK,
			Data:       reportData,
		})
	})

	// load robots list
	// swagger:route Post /robots/load Robots idOfRobotsLoad
	// Percentage
	//
	// responses:
	//   200:
	v1.Post("/robots/load", func(c *fiber.Ctx) error {
		resp, err := http.Get("https://robotstakeover20210903110417.azurewebsites.net/robotcpu")
		if err != nil {
			log.Fatalln(err)
		}
		//We Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error("unable to fetch the robots list", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		var robotList []models.RobotList
		err = json.Unmarshal(body, &robotList)

		if err != nil {
			logger.Error("unable to fetch the robots list", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}

		// load data to table
		err = handler.LoadRobotsHandler(robotList)
		if err != nil {
			logger.Error("unable to store the robots list", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		return c.Status(http.StatusOK).JSON(models.APIResponse{
			StatusCode: http.StatusOK,
			Message:    "successfully loaded robot list",
		})
	})

	// load robots list
	// swagger:route Post /robots/load Robots idOfRobotsLoad
	// Percentage
	//
	// responses:
	//   200:
	v1.Get("/robots/list", func(c *fiber.Ctx) error {

		data, err := handler.ListRobotsHandler()
		if err != nil {
			logger.Error("unable to store the robots list", zap.Error(err))
			return c.Status(http.StatusBadRequest).JSON(models.APIResponse{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			})
		}
		return c.Status(http.StatusOK).JSON(models.APIResponse{
			StatusCode: http.StatusOK,
			Message:    "successfully loaded robot list",
			Data:       data,
		})
	})
}
