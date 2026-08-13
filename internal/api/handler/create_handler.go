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

// Handle godoc
// @Summary      Create a new script
// @Description  Creates a new bash script on a remote host via SSH
// @Tags         scripts
// @Accept       json
// @Produce      json
// @Param        request body domain.CreateScriptRequest true "Create script request"
// @Success      201  {object}  domain.CreateScriptResponse
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /create [post]
func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateScriptRequest
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

	script, err := h.scriptService.CreateScript(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError

		switch {
		case errors.Is(err, domain.ErrTemplateNotFound):
			status = http.StatusBadRequest
		case errors.Is(err, domain.ErrScriptAlreadyExists):
			status = http.StatusConflict
		case errors.Is(err, domain.ErrInvalidHost), errors.Is(err, domain.ErrInvalidUser):
			status = http.StatusBadRequest
		}

		response.Error(w, status, err.Error())
		return
	}

	resp := domain.CreateScriptResponse{
		ScriptID:   script.ID.String(),
		ScriptPath: script.Path,
		CreatedAt:  script.CreatedAt,
	}

	response.JSON(w, http.StatusCreated, resp)
}
