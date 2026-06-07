package dtos

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"f1-race-events/models"
)

type CreateRaceEventRequest struct {
	RaceID      string `json:"raceId"`
	DriverID    string `json:"driverId"`
	ScuderiaID  string `json:"scuderiaId"`
	Lap         int    `json:"lap"`
	EventType   string `json:"eventType"`
	Description string `json:"description"`
}

type OutputRaceEventRequest struct {
	EventID     string `json:"eventId"`
	RaceID      string `json:"raceId"`
	DriverID    string `json:"driverId"`
	ScuderiaID  string `json:"scuderiaId"`
	Lap         int    `json:"lap"`
	EventType   string `json:"eventType"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
}

type OutputRacesEventsRequest struct {
	Events []OutputRaceEventRequest `json:"events"`
}

type GetRaceEventsQuery struct {
	RaceID        string
	DriverID      string
	ScuderiaID    string
	EventType     string
	Lap           int
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
}

func NewGetRaceEventsQuery(values url.Values) (GetRaceEventsQuery, error) {
	query := GetRaceEventsQuery{
		RaceID:     strings.TrimSpace(values.Get("raceId")),
		DriverID:   strings.TrimSpace(values.Get("driverId")),
		ScuderiaID: strings.TrimSpace(values.Get("scuderiaId")),
		EventType:  strings.TrimSpace(values.Get("eventType")),
	}

	lapParam := strings.TrimSpace(values.Get("lap"))
	if lapParam != "" {
		parsedLap, err := strconv.Atoi(lapParam)
		if err != nil {
			return query, errors.New("lap must be a number")
		}
		if parsedLap <= 0 {
			return query, errors.New("lap must be greater than 0")
		}
		query.Lap = parsedLap
	}

	createdAtParam := strings.TrimSpace(values.Get("createdAt"))
	createdAtFromParam := strings.TrimSpace(values.Get("createdAtFrom"))
	createdAtToParam := strings.TrimSpace(values.Get("createdAtTo"))

	if createdAtParam != "" && (createdAtFromParam != "" || createdAtToParam != "") {
		return query, errors.New("createdAt cannot be combined with createdAtFrom or createdAtTo")
	}

	if createdAtParam != "" {
		parsed, err := parseCreatedAt(createdAtParam)
		if err != nil {
			return query, errors.New("createdAt must be in format 2006-01-02T15:04:05 or RFC3339")
		}
		start := parsed.UTC()
		end := start.Add(time.Second - time.Nanosecond)
		query.CreatedAtFrom = &start
		query.CreatedAtTo = &end
	}

	if createdAtFromParam != "" {
		parsedFrom, err := parseCreatedAt(createdAtFromParam)
		if err != nil {
			return query, errors.New("createdAtFrom must be in format 2006-01-02T15:04:05 or RFC3339")
		}
		parsedFrom = parsedFrom.UTC()
		query.CreatedAtFrom = &parsedFrom
	}

	if createdAtToParam != "" {
		parsedTo, err := parseCreatedAt(createdAtToParam)
		if err != nil {
			return query, errors.New("createdAtTo must be in format 2006-01-02T15:04:05 or RFC3339")
		}
		parsedTo = parsedTo.UTC()
		query.CreatedAtTo = &parsedTo
	}

	if query.CreatedAtFrom != nil && query.CreatedAtTo != nil && query.CreatedAtFrom.After(*query.CreatedAtTo) {
		return query, errors.New("createdAtFrom must be before or equal to createdAtTo")
	}

	if strings.TrimSpace(query.RaceID) == "" {
		return query, errors.New("raceId is required")
	}

	if query.EventType != "" && !models.EventType(query.EventType).IsValid() {
		return query, errors.New("invalid eventType")
	}

	return query, nil
}

func parseCreatedAt(value string) (time.Time, error) {
	layouts := []string{time.RFC3339, "2006-01-02T15:04:05"}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, errors.New("invalid time format")
}

func (r CreateRaceEventRequest) Validate() error {
	var errs []string

	if strings.TrimSpace(r.RaceID) == "" {
		errs = append(errs, "raceId is required")
	}

	if strings.TrimSpace(r.DriverID) == "" {
		errs = append(errs, "driverId is required")
	}

	if strings.TrimSpace(r.ScuderiaID) == "" {
		errs = append(errs, "scuderiaId is required")
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
