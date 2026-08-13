package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/service"
	"github.com/go-playground/validator/v10"
)

type CallbackHandler struct {
	eventService *service.EventService
	validator    *validator.Validate
}

func NewCallbackHandler(eventService *service.EventService) *CallbackHandler {
	return &CallbackHandler{
		eventService: eventService,
		validator:    validator.New(),
	}
}

func (h *CallbackHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req domain.CallbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.eventService.ProcessCallback(r.Context(), req); err != nil {
		status := http.StatusInternalServerError
		if err == domain.ErrScriptNotFound {
			status = http.StatusNotFound
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, domain.CallbackResponse{
		Status: "ok",
	})
}
