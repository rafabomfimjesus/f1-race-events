package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"f1-race-events/dtos"
	"f1-race-events/usecases"
)

func RegisterRoutes(
	createUsecase *usecases.CreateRaceEventUseCase,
	listUsecase *usecases.GetAllRaceEventsUseCase,
) {
	http.HandleFunc("/race-events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateRaceEventHandler(createUsecase, w, r)
		case http.MethodGet:
			ListRaceEventsHandler(listUsecase, w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
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
		EventID:     evt.EventID,
		RaceID:      evt.RaceID,
		DriverID:    evt.DriverID,
		ScuderiaID:  evt.ScuderiaID,
		Lap:         evt.Lap,
		EventType:   string(evt.EventType),
		Description: evt.Description,
		CreatedAt:   evt.CreatedAt.Format(time.RFC3339),
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(response).Encode(out)
}

func ListRaceEventsHandler(
	usecase *usecases.GetAllRaceEventsUseCase,
	response http.ResponseWriter,
	router *http.Request,
) {
	query, err := dtos.NewGetRaceEventsQuery(router.URL.Query())
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	filters := usecases.GetAllRaceEventsFilters{
		RaceID:        query.RaceID,
		DriverID:      query.DriverID,
		ScuderiaID:    query.ScuderiaID,
		EventType:     query.EventType,
		Lap:           query.Lap,
		CreatedAtFrom: query.CreatedAtFrom,
		CreatedAtTo:   query.CreatedAtTo,
	}

	result, err := usecase.Execute(filters)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(response).Encode(result)
}
