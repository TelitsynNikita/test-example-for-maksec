package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/api/response"
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

// Handle godoc
// @Summary      Receive callback from agent
// @Description  Receives execution events from monitoring agent
// @Tags         callbacks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CallbackRequest true "Callback request"
// @Success      200  {object}  domain.CallbackResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /callback [post]
func (h *CallbackHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req domain.CallbackRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.eventService.ProcessCallback(r.Context(), req); err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, domain.ErrScriptNotFound) {
			status = http.StatusNotFound
		}

		response.Error(w, status, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, domain.CallbackResponse{
		Status: "ok",
	})
}
