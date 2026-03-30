package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kubedeck/kubedeck/internal/controllers"
	gitpkg "github.com/kubedeck/kubedeck/internal/git"
	"github.com/kubedeck/kubedeck/internal/models"
)

type WebhookHandler struct {
	apps      *models.AppStore
	buildCtrl *controllers.BuildController
}

func NewWebhookHandler(apps *models.AppStore, buildCtrl *controllers.BuildController) *WebhookHandler {
	return &WebhookHandler{apps: apps, buildCtrl: buildCtrl}
}

func (h *WebhookHandler) GitHub(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "app_id")
	app, err := h.apps.GetByID(appID)
	if err != nil {
		writeError(w, http.StatusNotFound, "app not found")
		return
	}

	payload, err := gitpkg.ParseGitHubWebhook(r, app.WebhookSecret)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid webhook: "+err.Error())
		return
	}

	if payload.Branch != app.GitBranch {
		writeJSON(w, http.StatusOK, map[string]string{"message": "branch mismatch, skipping"})
		return
	}

	if !app.AutoDeploy {
		writeJSON(w, http.StatusOK, map[string]string{"message": "auto deploy disabled, skipping"})
		return
	}

	build, err := h.buildCtrl.TriggerBuild(app, payload.CommitSHA, payload.CommitMessage, payload.CommitAuthor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to trigger build")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "build triggered", "build_id": build.ID})
}

func (h *WebhookHandler) GitLab(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "app_id")
	app, err := h.apps.GetByID(appID)
	if err != nil {
		writeError(w, http.StatusNotFound, "app not found")
		return
	}

	payload, err := gitpkg.ParseGitLabWebhook(r, app.WebhookSecret)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid webhook: "+err.Error())
		return
	}

	if payload.Branch != app.GitBranch {
		writeJSON(w, http.StatusOK, map[string]string{"message": "branch mismatch, skipping"})
		return
	}

	if !app.AutoDeploy {
		writeJSON(w, http.StatusOK, map[string]string{"message": "auto deploy disabled, skipping"})
		return
	}

	build, err := h.buildCtrl.TriggerBuild(app, payload.CommitSHA, payload.CommitMessage, payload.CommitAuthor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to trigger build")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "build triggered", "build_id": build.ID})
}
