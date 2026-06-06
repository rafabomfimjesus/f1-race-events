package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"f1-race-events/dtos"
	"f1-race-events/usecases"
)

func RegisterRoutes(usecase *usecases.CreateRaceEventUseCase) {
	http.HandleFunc("/race-events", func(w http.ResponseWriter, r *http.Request) {
		CreateRaceEventHandler(usecase, w, r)
	})
}

func CreateRaceEventHandler(
	usecase *usecases.CreateRaceEventUseCase,
	response http.ResponseWriter,
	router *http.Request,
) {
	if router.Method != http.MethodPost {
		response.Header().Set("Allow", http.MethodPost)
		http.Error(response, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request dtos.CreateRaceEventRequest
	if err := json.NewDecoder(router.Body).Decode(&request); err != nil {
		http.Error(response, "invalid request body", http.StatusBadRequest)
		return
	}
	evt, err := usecase.Execute(request)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	out := dtos.OutputRaceEventRequest{
		RaceID:      evt.RaceID,
		DriverID:    evt.DriverID,
		Lap:         evt.Lap,
		EventType:   string(evt.EventType),
		Description: evt.Description,
		CreatedAt:   evt.CreatedAt.Format(time.RFC3339),
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(response).Encode(out)
}
