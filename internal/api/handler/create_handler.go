package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/domain"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/service"
	"github.com/go-playground/validator/v10"
)

type CreateHandler struct {
	scriptService *service.ScriptService
	validator     *validator.Validate
}

func NewCreateHandler(scriptService *service.ScriptService) *CreateHandler {
	return &CreateHandler{
		scriptService: scriptService,
		validator:     validator.New(),
	}
}

func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateScriptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	script, err := h.scriptService.CreateScript(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case domain.ErrTemplateNotFound:
			status = http.StatusBadRequest
		case domain.ErrScriptAlreadyExists:
			status = http.StatusConflict
		case domain.ErrInvalidHost, domain.ErrInvalidUser:
			status = http.StatusBadRequest
		}
		respondWithError(w, status, err.Error())
		return
	}

	resp := domain.CreateScriptResponse{
		ScriptID:   script.ID.String(),
		ScriptPath: script.Path,
		CreatedAt:  script.CreatedAt,
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}
