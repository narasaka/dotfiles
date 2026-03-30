package api

import (
	"encoding/json"
	"net/http"

	"github.com/kubeploy/kubeploy/internal/models"
)

type SettingsHandler struct {
	settings *models.SettingsStore
}

func NewSettingsHandler(settings *models.SettingsStore) *SettingsHandler {
	return &SettingsHandler{settings: settings}
}

func (h *SettingsHandler) Get(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settings.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get settings")
		return
	}

	result := make(map[string]string)
	for _, s := range settings {
		// Don't expose sensitive values
		if s.Key == "registry_password" || s.Key == "session_secret" {
			result[s.Key] = "••••••••"
		} else {
			result[s.Key] = s.Value
		}
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	var updates map[string]string
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Don't allow updating session_secret via API
	delete(updates, "session_secret")

	// Skip masked values
	for key, value := range updates {
		if value == "••••••••" {
			delete(updates, key)
		}
	}

	if err := h.settings.SetMultiple(updates); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update settings")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "settings updated"})
}
