package main

import (
	"net/http"

	"f1-race-events/repositories"
	"f1-race-events/routes"
	"f1-race-events/usecases"
	"f1-race-events/utils"
)

func main() {

	logger := utils.New()

	logger.Info("starting application")

	client, err := utils.NewDynamoClient()

	if err != nil {
		logger.Error(
			"failed to create dynamodb client",
			"error",
			err,
		)
		return
	}

	repository := repositories.InitRaceEventRepository(
		client,
		"race_event",
	)

	usecase := usecases.NewCreateRaceEventUseCase(repository)

	routes.RegisterRoutes(usecase)

	logger.Info(
		"server listening",
		"port",
		8080,
	)

	err = http.ListenAndServe(
		":8080",
		nil,
	)

	if err != nil {
		logger.Error(
			"failed to start server",
			"error",
			err,
		)
	}
}
