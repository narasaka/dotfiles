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

	kubedeck "github.com/kubedeck/kubedeck"
	"github.com/kubedeck/kubedeck/internal/api"
	"github.com/kubedeck/kubedeck/internal/config"
	"github.com/kubedeck/kubedeck/internal/controllers"
	"github.com/kubedeck/kubedeck/internal/db"
	"github.com/kubedeck/kubedeck/internal/k8s"
	"github.com/kubedeck/kubedeck/internal/models"
)

func main() {
	cfg := config.Load()

	log.Printf("Starting Kubedeck on port %d (dev=%v)", cfg.Port, cfg.Dev)

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
		webFS, err = fs.Sub(kubedeck.WebDist, "web/dist")
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
		log.Printf("Kubedeck is running at http://localhost:%d", cfg.Port)
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
