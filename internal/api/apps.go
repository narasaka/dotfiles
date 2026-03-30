package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kubeploy/kubeploy/internal/models"
)

type AppHandler struct {
	apps   *models.AppStore
	builds *models.BuildStore
}

func NewAppHandler(apps *models.AppStore, builds *models.BuildStore) *AppHandler {
	return &AppHandler{apps: apps, builds: builds}
}

func (h *AppHandler) List(w http.ResponseWriter, r *http.Request) {
	apps, err := h.apps.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list apps")
		return
	}
	writeJSON(w, http.StatusOK, apps)
}

func (h *AppHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	app, err := h.apps.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "app not found")
		return
	}
	writeJSON(w, http.StatusOK, app)
}

type createAppRequest struct {
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	GitURL         string `json:"git_url"`
	GitBranch      string `json:"git_branch"`
	GitSubpath     string `json:"git_subpath"`
	DockerfilePath string `json:"dockerfile_path"`
	RegistryImage  string `json:"registry_image"`
	Namespace      string `json:"namespace"`
	Replicas       int    `json:"replicas"`
	Port           int    `json:"port"`
	EnvVars        string `json:"env_vars"`
	AutoDeploy     bool   `json:"auto_deploy"`
	IngressHost    string `json:"ingress_host"`
	IngressTLS     bool   `json:"ingress_tls"`
}

func (h *AppHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.GitURL == "" {
		writeError(w, http.StatusBadRequest, "name and git_url are required")
		return
	}

	if req.GitBranch == "" {
		req.GitBranch = "main"
	}
	if req.DockerfilePath == "" {
		req.DockerfilePath = "Dockerfile"
	}
	if req.Namespace == "" {
		req.Namespace = "default"
	}
	if req.Replicas == 0 {
		req.Replicas = 1
	}
	if req.Port == 0 {
		req.Port = 8080
	}
	if req.EnvVars == "" {
		req.EnvVars = "{}"
	}

	app := &models.App{
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		GitURL:         req.GitURL,
		GitBranch:      req.GitBranch,
		GitSubpath:     req.GitSubpath,
		DockerfilePath: req.DockerfilePath,
		RegistryImage:  req.RegistryImage,
		Namespace:      req.Namespace,
		Replicas:       req.Replicas,
		Port:           req.Port,
		EnvVars:        req.EnvVars,
		AutoDeploy:     req.AutoDeploy,
		IngressHost:    req.IngressHost,
		IngressTLS:     req.IngressTLS,
	}

	created, err := h.apps.Create(app)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create app")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *AppHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	existing, err := h.apps.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "app not found")
		return
	}

	var req createAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	existing.Name = req.Name
	existing.DisplayName = req.DisplayName
	existing.GitURL = req.GitURL
	existing.GitBranch = req.GitBranch
	existing.GitSubpath = req.GitSubpath
	existing.DockerfilePath = req.DockerfilePath
	existing.RegistryImage = req.RegistryImage
	existing.Namespace = req.Namespace
	existing.Replicas = req.Replicas
	existing.Port = req.Port
	existing.AutoDeploy = req.AutoDeploy
	existing.IngressHost = req.IngressHost
	existing.IngressTLS = req.IngressTLS

	if req.EnvVars != "" {
		existing.EnvVars = req.EnvVars
	}

	updated, err := h.apps.Update(existing)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update app")
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *AppHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.apps.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete app")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func (h *AppHandler) GetEnv(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	app, err := h.apps.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "app not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(app.EnvVars))
}

type updateEnvRequest struct {
	EnvVars string `json:"env_vars"`
}

func (h *AppHandler) UpdateEnv(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req updateEnvRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.apps.UpdateEnvVars(id, req.EnvVars); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update env vars")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *AppHandler) Redeploy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	app, err := h.apps.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "app not found")
		return
	}

	if app.CurrentBuildID == nil {
		writeError(w, http.StatusBadRequest, "no build to redeploy")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "redeploy triggered", "app_id": id})
}
