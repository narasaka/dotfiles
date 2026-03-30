package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kubeploy/kubeploy/internal/controllers"
	"github.com/kubeploy/kubeploy/internal/models"
)

type DeploymentHandler struct {
	deployments *models.DeploymentStore
	apps        *models.AppStore
	builds      *models.BuildStore
	deployCtrl  *controllers.DeployController
}

func NewDeploymentHandler(
	deployments *models.DeploymentStore,
	apps *models.AppStore,
	builds *models.BuildStore,
	deployCtrl *controllers.DeployController,
) *DeploymentHandler {
	return &DeploymentHandler{
		deployments: deployments,
		apps:        apps,
		builds:      builds,
		deployCtrl:  deployCtrl,
	}
}

func (h *DeploymentHandler) ListByApp(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "id")
	deployments, err := h.deployments.ListByApp(appID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list deployments")
		return
	}
	writeJSON(w, http.StatusOK, deployments)
}

func (h *DeploymentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	dep, err := h.deployments.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "deployment not found")
		return
	}
	writeJSON(w, http.StatusOK, dep)
}

func (h *DeploymentHandler) Rollback(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	dep, err := h.deployments.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "deployment not found")
		return
	}

	app, err := h.apps.GetByID(dep.AppID)
	if err != nil {
		writeError(w, http.StatusNotFound, "app not found")
		return
	}

	build, err := h.builds.GetByID(dep.BuildID)
	if err != nil {
		writeError(w, http.StatusNotFound, "build not found")
		return
	}

	newDep, err := h.deployCtrl.Deploy(app, build)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to rollback: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, newDep)
}
