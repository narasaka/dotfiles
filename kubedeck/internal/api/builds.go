package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kubedeck/kubedeck/internal/controllers"
	"github.com/kubedeck/kubedeck/internal/models"
)

type BuildHandler struct {
	builds    *models.BuildStore
	apps      *models.AppStore
	buildCtrl *controllers.BuildController
}

func NewBuildHandler(builds *models.BuildStore, apps *models.AppStore, buildCtrl *controllers.BuildController) *BuildHandler {
	return &BuildHandler{builds: builds, apps: apps, buildCtrl: buildCtrl}
}

func (h *BuildHandler) ListByApp(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "id")
	builds, err := h.builds.ListByApp(appID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list builds")
		return
	}
	writeJSON(w, http.StatusOK, builds)
}

func (h *BuildHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	build, err := h.builds.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "build not found")
		return
	}
	writeJSON(w, http.StatusOK, build)
}

type triggerBuildRequest struct {
	CommitSHA     string `json:"commit_sha"`
	CommitMessage string `json:"commit_message"`
	CommitAuthor  string `json:"commit_author"`
}

func (h *BuildHandler) TriggerBuild(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "id")
	app, err := h.apps.GetByID(appID)
	if err != nil {
		writeError(w, http.StatusNotFound, "app not found")
		return
	}

	var req triggerBuildRequest
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&req)
	}

	if req.CommitSHA == "" {
		req.CommitSHA = "manual"
	}

	build, err := h.buildCtrl.TriggerBuild(app, req.CommitSHA, req.CommitMessage, req.CommitAuthor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to trigger build: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, build)
}

func (h *BuildHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	build, err := h.builds.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "build not found")
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(build.Logs))
}

func (h *BuildHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	build, err := h.builds.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "build not found")
		return
	}

	if build.Status != "pending" && build.Status != "building" {
		writeError(w, http.StatusBadRequest, "build is not cancellable")
		return
	}

	if err := h.buildCtrl.CancelBuild(build); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to cancel build")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "build cancelled"})
}
