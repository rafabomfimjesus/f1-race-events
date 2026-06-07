package usecases

import (
	"context"
	"time"

	"f1-race-events/dtos"
	"f1-race-events/repositories"
)

type GetAllRaceEventsFilters struct {
	RaceID        string
	DriverID      string
	ScuderiaID    string
	EventType     string
	Lap           int
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
}

type GetAllRaceEventsUseCase struct {
	repository *repositories.RaceEventRepository
}

func NewGetAllRaceEventsUseCase(
	repository *repositories.RaceEventRepository,
) *GetAllRaceEventsUseCase {
	return &GetAllRaceEventsUseCase{
		repository: repository,
	}
}

func (u *GetAllRaceEventsUseCase) Execute(
	filters GetAllRaceEventsFilters,
) (dtos.OutputRacesEventsRequest, error) {
	queryFilters := repositories.RaceEventQueryFilters{
		RaceID:        "RACE#" + filters.RaceID,
		DriverID:      filters.DriverID,
		ScuderiaID:    filters.ScuderiaID,
		EventType:     filters.EventType,
		Lap:           filters.Lap,
		CreatedAtFrom: filters.CreatedAtFrom,
		CreatedAtTo:   filters.CreatedAtTo,
	}

	events, err := u.repository.GetByRaceID(context.Background(), queryFilters)
	if err != nil {
		return dtos.OutputRacesEventsRequest{}, err
	}

	outputEvents := make([]dtos.OutputRaceEventRequest, 0, len(events))
	for _, evt := range events {
		outputEvents = append(outputEvents, dtos.OutputRaceEventRequest{
			EventID:     evt.EventID,
			RaceID:      evt.RaceID,
			DriverID:    evt.DriverID,
			ScuderiaID:  evt.ScuderiaID,
			Lap:         evt.Lap,
			EventType:   string(evt.EventType),
			Description: evt.Description,
			CreatedAt:   evt.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return dtos.OutputRacesEventsRequest{
		Events: outputEvents,
	}, nil
}
