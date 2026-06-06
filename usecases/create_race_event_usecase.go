package usecases

import (
	"context"
	"time"

	"f1-race-events/dtos"
	"f1-race-events/models"
	"f1-race-events/repositories"

	"github.com/google/uuid"
)

type CreateRaceEventUseCase struct {
	repository *repositories.RaceEventRepository
}

func NewCreateRaceEventUseCase(
	repository *repositories.RaceEventRepository,
) *CreateRaceEventUseCase {
	return &CreateRaceEventUseCase{
		repository: repository,
	}
}

func (u *CreateRaceEventUseCase) Execute(
	request dtos.CreateRaceEventRequest,
) (models.RaceEvent, error) {

	if err := request.Validate(); err != nil {
		return models.RaceEvent{}, err
	}

	eventID := uuid.NewString()
	raceID := request.RaceID
	driverID := request.DriverID
	eventType := models.EventType(request.EventType)
	description := request.Description
	createdAt := time.Now().UTC()

	event := models.RaceEvent{
		EventID:     eventID,
		RaceID:      raceID,
		DriverID:    driverID,
		Lap:         request.Lap,
		EventType:   eventType,
		Description: description,
		CreatedAt:   createdAt,
	}

	if err := u.repository.Create(context.Background(), event); err != nil {
		return models.RaceEvent{}, err
	}

	return event, nil
}
