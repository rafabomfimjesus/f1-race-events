package dtos

import (
	"errors"
	"strings"

	"f1-race-events/models"
)

type CreateRaceEventRequest struct {
	RaceID      string `json:"raceId"`
	DriverID    string `json:"driverId"`
	Lap         int    `json:"lap"`
	EventType   string `json:"eventType"`
	Description string `json:"description"`
}

type OutputRaceEventRequest struct {
	RaceID      string `json:"raceId"`
	DriverID    string `json:"driverId"`
	Lap         int    `json:"lap"`
	EventType   string `json:"eventType"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

func (r CreateRaceEventRequest) Validate() error {
	var errs []string

	if strings.TrimSpace(r.RaceID) == "" {
		errs = append(errs, "raceId is required")
	}

	if strings.TrimSpace(r.DriverID) == "" {
		errs = append(errs, "driverId is required")
	}

	if r.Lap == 0 {
		errs = append(errs, "lap must be greater than 0")
	}

	if strings.TrimSpace(r.EventType) == "" {
		errs = append(errs, "eventType is required")
	}

	if !models.EventType(r.EventType).IsValid() {
		errs = append(errs, "invalid eventType")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}

	return nil
}
