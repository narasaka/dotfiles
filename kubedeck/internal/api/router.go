package api

import (
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/kubedeck/kubedeck/internal/controllers"
	"github.com/kubedeck/kubedeck/internal/k8s"
	"github.com/kubedeck/kubedeck/internal/models"
)

type Server struct {
	router      chi.Router
	auth        *AuthMiddleware
	authHandler *AuthHandler
	appHandler  *AppHandler
	buildHandler *BuildHandler
	deployHandler *DeploymentHandler
	logHandler   *LogHandler
	webhookHandler *WebhookHandler
	settingsHandler *SettingsHandler
	webFS       fs.FS
}

type ServerDeps struct {
	Users       *models.UserStore
	Apps        *models.AppStore
	Builds      *models.BuildStore
	Deployments *models.DeploymentStore
	Settings    *models.SettingsStore
	BuildCtrl   *controllers.BuildController
	DeployCtrl  *controllers.DeployController
	K8sClient   *k8s.Client
	SessionSecret string
	DevMode     bool
	WebFS       fs.FS
}

func NewServer(deps ServerDeps) *Server {
	authMW := NewAuthMiddleware(deps.Users, deps.SessionSecret)

	s := &Server{
		auth:            authMW,
		authHandler:     NewAuthHandler(deps.Users, authMW),
		appHandler:      NewAppHandler(deps.Apps, deps.Builds),
		buildHandler:    NewBuildHandler(deps.Builds, deps.Apps, deps.BuildCtrl),
		deployHandler:   NewDeploymentHandler(deps.Deployments, deps.Apps, deps.Builds, deps.DeployCtrl),
		logHandler:      NewLogHandler(deps.Apps, deps.Builds, deps.K8sClient),
		webhookHandler:  NewWebhookHandler(deps.Apps, deps.BuildCtrl),
		settingsHandler: NewSettingsHandler(deps.Settings),
		webFS:           deps.WebFS,
	}

	r := chi.NewRouter()

	// Global middleware
	r.Use(RecoveryMiddleware)
	r.Use(LoggingMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/auth/setup", s.authHandler.Setup)
		r.Post("/auth/login", s.authHandler.Login)
		r.Get("/auth/check", s.authHandler.CheckSetup)

		// Webhooks (no auth — validated by HMAC)
		r.Post("/webhooks/github/{app_id}", s.webhookHandler.GitHub)
		r.Post("/webhooks/gitlab/{app_id}", s.webhookHandler.GitLab)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(s.auth.Require)

			// Auth
			r.Post("/auth/logout", s.authHandler.Logout)
			r.Get("/auth/me", s.authHandler.Me)

			// Apps
			r.Get("/apps", s.appHandler.List)
			r.Post("/apps", s.appHandler.Create)
			r.Get("/apps/{id}", s.appHandler.Get)
			r.Put("/apps/{id}", s.appHandler.Update)
			r.Delete("/apps/{id}", s.appHandler.Delete)
			r.Get("/apps/{id}/env", s.appHandler.GetEnv)
			r.Put("/apps/{id}/env", s.appHandler.UpdateEnv)
			r.Post("/apps/{id}/redeploy", s.appHandler.Redeploy)

			// Builds
			r.Get("/apps/{id}/builds", s.buildHandler.ListByApp)
			r.Post("/apps/{id}/builds", s.buildHandler.TriggerBuild)
			r.Get("/builds/{id}", s.buildHandler.Get)
			r.Get("/builds/{id}/logs", s.buildHandler.GetLogs)
			r.Post("/builds/{id}/cancel", s.buildHandler.Cancel)

			// Deployments
			r.Get("/apps/{id}/deployments", s.deployHandler.ListByApp)
			r.Get("/deployments/{id}", s.deployHandler.Get)
			r.Post("/deployments/{id}/rollback", s.deployHandler.Rollback)

			// Settings
			r.Get("/settings", s.settingsHandler.Get)
			r.Put("/settings", s.settingsHandler.Update)
		})

		// WebSocket routes (auth checked in handler)
		r.Get("/apps/{id}/logs/ws", s.logHandler.StreamAppLogs)
		r.Get("/builds/{id}/logs/ws", s.logHandler.StreamBuildLogs)
	})

	// Serve frontend
	if s.webFS != nil {
		fileServer := http.FileServer(http.FS(s.webFS))
		r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to serve the file directly
			f, err := s.webFS.Open(r.URL.Path[1:])
			if err != nil {
				// Fall back to index.html for SPA routing
				r.URL.Path = "/"
			} else {
				f.Close()
			}
			fileServer.ServeHTTP(w, r)
		}))
	}

	s.router = r
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
