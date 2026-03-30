package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	kubeploy "github.com/kubeploy/kubeploy"
	"github.com/kubeploy/kubeploy/internal/api"
	"github.com/kubeploy/kubeploy/internal/config"
	"github.com/kubeploy/kubeploy/internal/controllers"
	"github.com/kubeploy/kubeploy/internal/db"
	"github.com/kubeploy/kubeploy/internal/k8s"
	"github.com/kubeploy/kubeploy/internal/models"

	// Register provider plugins
	_ "github.com/kubeploy/kubeploy/internal/plugins/gke"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting Kubeploy on port %d (dev=%v)", cfg.Port, cfg.Dev)

	// Open database
	database, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// Initialize stores
	userStore := models.NewUserStore(database)
	appStore := models.NewAppStore(database)
	buildStore := models.NewBuildStore(database)
	deploymentStore := models.NewDeploymentStore(database)
	settingsStore := models.NewSettingsStore(database)
	clusterStore := models.NewClusterStore(database)

	// Initialize K8s client (optional in dev mode)
	var k8sClient *k8s.Client
	if !cfg.Dev {
		k8sClient, err = k8s.NewClient(cfg.Namespace)
		if err != nil {
			log.Printf("Warning: K8s client not available: %v", err)
		}
	}

	// Initialize controllers
	deployCtrl := controllers.NewDeployController(deploymentStore, appStore, k8sClient)
	buildCtrl := controllers.NewBuildController(buildStore, appStore, settingsStore, k8sClient, deployCtrl)

	// Start cleanup controller
	cleanupCtrl := controllers.NewCleanupController(k8sClient)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go cleanupCtrl.Start(ctx)

	// Resolve session secret
	sessionSecret := cfg.SessionSecret
	if sessionSecret == "" {
		sessionSecret, _ = settingsStore.Get("session_secret")
	}

	// Prepare web filesystem
	var webFS fs.FS
	if !cfg.Dev {
		webFS, err = fs.Sub(kubeploy.WebDist, "web/dist")
		if err != nil {
			log.Printf("Warning: embedded web assets not available: %v", err)
		}
	}

	// Create API server
	server := api.NewServer(api.ServerDeps{
		Users:         userStore,
		Apps:          appStore,
		Builds:        buildStore,
		Deployments:   deploymentStore,
		Settings:      settingsStore,
		Clusters:      clusterStore,
		BuildCtrl:     buildCtrl,
		DeployCtrl:    deployCtrl,
		K8sClient:     k8sClient,
		SessionSecret: sessionSecret,
		DevMode:       cfg.Dev,
		WebFS:         webFS,
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: server,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Kubeploy is running at http://localhost:%d", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	httpServer.Shutdown(shutdownCtx)
}
