package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"nhooyr.io/websocket"

	"github.com/kubeploy/kubeploy/internal/k8s"
	"github.com/kubeploy/kubeploy/internal/models"
)

type LogHandler struct {
	apps      *models.AppStore
	builds    *models.BuildStore
	k8sClient *k8s.Client
}

func NewLogHandler(apps *models.AppStore, builds *models.BuildStore, k8sClient *k8s.Client) *LogHandler {
	return &LogHandler{apps: apps, builds: builds, k8sClient: k8sClient}
}

func (h *LogHandler) StreamAppLogs(w http.ResponseWriter, r *http.Request) {
	appID := chi.URLParam(r, "id")
	app, err := h.apps.GetByID(appID)
	if err != nil {
		writeError(w, http.StatusNotFound, "app not found")
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Printf("websocket accept error: %v", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	labelSelector := fmt.Sprintf("kubeploy/app-id=%s", app.ID)

	if h.k8sClient == nil {
		conn.Write(ctx, websocket.MessageText, []byte("K8s client not available (dev mode)\n"))
		<-ctx.Done()
		return
	}

	logCh, err := h.k8sClient.StreamPodLogs(ctx, app.Namespace, labelSelector)
	if err != nil {
		conn.Write(ctx, websocket.MessageText, []byte("Error streaming logs: "+err.Error()+"\n"))
		return
	}

	for line := range logCh {
		if err := conn.Write(ctx, websocket.MessageText, []byte(line)); err != nil {
			return
		}
	}
}

func (h *LogHandler) StreamBuildLogs(w http.ResponseWriter, r *http.Request) {
	buildID := chi.URLParam(r, "id")
	build, err := h.builds.GetByID(buildID)
	if err != nil {
		writeError(w, http.StatusNotFound, "build not found")
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Printf("websocket accept error: %v", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Send existing logs first
	if build.Logs != "" {
		conn.Write(ctx, websocket.MessageText, []byte(build.Logs))
	}

	// If build is done, close
	if build.Status == "success" || build.Status == "failed" || build.Status == "cancelled" {
		return
	}

	// Poll for new logs
	lastLen := len(build.Logs)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			build, err = h.builds.GetByID(buildID)
			if err != nil {
				return
			}

			if len(build.Logs) > lastLen {
				newLogs := build.Logs[lastLen:]
				if err := conn.Write(ctx, websocket.MessageText, []byte(newLogs)); err != nil {
					return
				}
				lastLen = len(build.Logs)
			}

			if build.Status == "success" || build.Status == "failed" || build.Status == "cancelled" {
				return
			}
		}
	}
}
