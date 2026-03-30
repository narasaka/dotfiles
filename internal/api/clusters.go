package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/kubeploy/kubeploy/internal/models"
	"github.com/kubeploy/kubeploy/internal/plugins"
)

type ClusterHandler struct {
	clusters *models.ClusterStore
}

func NewClusterHandler(clusters *models.ClusterStore) *ClusterHandler {
	return &ClusterHandler{clusters: clusters}
}

// --- Provider endpoints ---

func (h *ClusterHandler) ListProviders(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, plugins.List())
}

func (h *ClusterHandler) ValidateCredentials(w http.ResponseWriter, r *http.Request) {
	providerName := chi.URLParam(r, "provider")
	provider, err := plugins.Get(providerName)
	if err != nil {
		writeError(w, http.StatusBadRequest, "unknown provider: "+providerName)
		return
	}

	var creds map[string]string
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := provider.ValidateCredentials(r.Context(), creds); err != nil {
		writeError(w, http.StatusBadRequest, "invalid credentials: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "credentials valid"})
}

func (h *ClusterHandler) DiscoverClusters(w http.ResponseWriter, r *http.Request) {
	providerName := chi.URLParam(r, "provider")
	provider, err := plugins.Get(providerName)
	if err != nil {
		writeError(w, http.StatusBadRequest, "unknown provider: "+providerName)
		return
	}

	var creds map[string]string
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	clusters, err := provider.ListClusters(r.Context(), creds)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list clusters: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, clusters)
}

// --- Cluster CRUD ---

func (h *ClusterHandler) List(w http.ResponseWriter, r *http.Request) {
	providerFilter := r.URL.Query().Get("provider")
	var clusters []models.Cluster
	var err error

	if providerFilter != "" {
		clusters, err = h.clusters.ListByProvider(providerFilter)
	} else {
		clusters, err = h.clusters.List()
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list clusters")
		return
	}
	writeJSON(w, http.StatusOK, clusters)
}

func (h *ClusterHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cluster, err := h.clusters.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "cluster not found")
		return
	}
	writeJSON(w, http.StatusOK, cluster)
}

type connectClusterRequest struct {
	Name              string            `json:"name"`
	DisplayName       string            `json:"display_name"`
	Provider          string            `json:"provider"`
	ProviderClusterID string            `json:"provider_cluster_id"`
	ProjectID         string            `json:"project_id"`
	Credentials       map[string]string `json:"credentials"`
}

func (h *ClusterHandler) Connect(w http.ResponseWriter, r *http.Request) {
	var req connectClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.Provider == "" || req.ProviderClusterID == "" {
		writeError(w, http.StatusBadRequest, "name, provider, and provider_cluster_id are required")
		return
	}

	provider, err := plugins.Get(req.Provider)
	if err != nil {
		writeError(w, http.StatusBadRequest, "unknown provider: "+req.Provider)
		return
	}

	// Validate credentials
	if err := provider.ValidateCredentials(r.Context(), req.Credentials); err != nil {
		writeError(w, http.StatusBadRequest, "invalid credentials: "+err.Error())
		return
	}

	// Fetch cluster details from provider
	detail, err := provider.GetCluster(r.Context(), req.Credentials, req.ProviderClusterID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to fetch cluster: "+err.Error())
		return
	}

	credsJSON, _ := json.Marshal(req.Credentials)

	cluster := &models.Cluster{
		Name:              req.Name,
		DisplayName:       req.DisplayName,
		Provider:          req.Provider,
		ProviderClusterID: req.ProviderClusterID,
		ProjectID:         req.ProjectID,
		Location:          detail.Location,
		Status:            "connected",
		K8sVersion:        detail.K8sVersion,
		Endpoint:          detail.Endpoint,
		NodeCount:         detail.NodeCount,
		Credentials:       string(credsJSON),
		Metadata:          "{}",
	}

	created, err := h.clusters.Create(cluster)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save cluster: "+err.Error())
		return
	}

	// Sync node pools
	h.syncNodePools(r.Context(), created, provider, req.Credentials)

	h.clusters.AddEvent(created.ID, "connected", "Cluster connected to Kubeploy", "{}")

	writeJSON(w, http.StatusCreated, created)
}

func (h *ClusterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.clusters.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete cluster")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "cluster disconnected"})
}

func (h *ClusterHandler) Sync(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cluster, err := h.clusters.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "cluster not found")
		return
	}

	provider, err := plugins.Get(cluster.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "provider not available: "+cluster.Provider)
		return
	}

	var creds map[string]string
	json.Unmarshal([]byte(cluster.Credentials), &creds)

	detail, err := provider.GetCluster(r.Context(), creds, cluster.ProviderClusterID)
	if err != nil {
		h.clusters.UpdateStatus(id, "error")
		h.clusters.AddEvent(id, "sync_error", "Sync failed: "+err.Error(), "{}")
		writeError(w, http.StatusInternalServerError, "failed to sync cluster: "+err.Error())
		return
	}

	cluster.K8sVersion = detail.K8sVersion
	cluster.Endpoint = detail.Endpoint
	cluster.NodeCount = detail.NodeCount
	cluster.Status = "connected"
	h.clusters.Update(cluster)

	h.syncNodePools(r.Context(), cluster, provider, creds)
	h.clusters.AddEvent(id, "sync", "Cluster synced successfully", "{}")

	writeJSON(w, http.StatusOK, cluster)
}

// --- Node Pool endpoints ---

func (h *ClusterHandler) ListNodePools(w http.ResponseWriter, r *http.Request) {
	clusterID := chi.URLParam(r, "id")
	pools, err := h.clusters.ListNodePools(clusterID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list node pools")
		return
	}
	writeJSON(w, http.StatusOK, pools)
}

type createNodePoolRequest struct {
	Name             string            `json:"name"`
	MachineType      string            `json:"machine_type"`
	DiskSizeGB       int               `json:"disk_size_gb"`
	DiskType         string            `json:"disk_type"`
	InitialNodeCount int               `json:"initial_node_count"`
	MinNodes         int               `json:"min_nodes"`
	MaxNodes         int               `json:"max_nodes"`
	Autoscaling      bool              `json:"autoscaling"`
	Preemptible      bool              `json:"preemptible"`
	SpotInstance     bool              `json:"spot_instance"`
	Labels           map[string]string `json:"labels"`
}

func (h *ClusterHandler) CreateNodePool(w http.ResponseWriter, r *http.Request) {
	clusterID := chi.URLParam(r, "id")
	cluster, err := h.clusters.GetByID(clusterID)
	if err != nil {
		writeError(w, http.StatusNotFound, "cluster not found")
		return
	}

	var req createNodePoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" || req.MachineType == "" {
		writeError(w, http.StatusBadRequest, "name and machine_type are required")
		return
	}

	provider, err := plugins.Get(cluster.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "provider not available")
		return
	}

	var creds map[string]string
	json.Unmarshal([]byte(cluster.Credentials), &creds)

	pool, err := provider.CreateNodePool(r.Context(), creds, cluster.ProviderClusterID, plugins.CreateNodePoolOpts{
		Name:             req.Name,
		MachineType:      req.MachineType,
		DiskSizeGB:       req.DiskSizeGB,
		DiskType:         req.DiskType,
		InitialNodeCount: req.InitialNodeCount,
		MinNodes:         req.MinNodes,
		MaxNodes:         req.MaxNodes,
		Autoscaling:      req.Autoscaling,
		Preemptible:      req.Preemptible,
		SpotInstance:     req.SpotInstance,
		Labels:           req.Labels,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create node pool: "+err.Error())
		return
	}

	labelsJSON, _ := json.Marshal(pool.Labels)
	taintsJSON, _ := json.Marshal(pool.Taints)

	dbPool := &models.ClusterNodePool{
		ClusterID:      clusterID,
		Name:           pool.Name,
		ProviderPoolID: pool.ID,
		MachineType:    pool.MachineType,
		DiskSizeGB:     pool.DiskSizeGB,
		DiskType:       pool.DiskType,
		NodeCount:      pool.NodeCount,
		MinNodes:       pool.MinNodes,
		MaxNodes:       pool.MaxNodes,
		Autoscaling:    pool.Autoscaling,
		Preemptible:    pool.Preemptible,
		SpotInstance:   pool.SpotInstance,
		Status:         pool.Status,
		K8sVersion:     pool.K8sVersion,
		Labels:         string(labelsJSON),
		Taints:         string(taintsJSON),
	}

	created, err := h.clusters.CreateNodePool(dbPool)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save node pool")
		return
	}

	h.clusters.AddEvent(clusterID, "node_pool_created", "Node pool '"+pool.Name+"' created", "{}")
	writeJSON(w, http.StatusCreated, created)
}

type scaleNodePoolRequest struct {
	NodeCount   *int  `json:"node_count,omitempty"`
	MinNodes    *int  `json:"min_nodes,omitempty"`
	MaxNodes    *int  `json:"max_nodes,omitempty"`
	Autoscaling *bool `json:"autoscaling,omitempty"`
}

func (h *ClusterHandler) ScaleNodePool(w http.ResponseWriter, r *http.Request) {
	clusterID := chi.URLParam(r, "id")
	poolID := chi.URLParam(r, "poolId")

	cluster, err := h.clusters.GetByID(clusterID)
	if err != nil {
		writeError(w, http.StatusNotFound, "cluster not found")
		return
	}

	pool, err := h.clusters.GetNodePool(poolID)
	if err != nil {
		writeError(w, http.StatusNotFound, "node pool not found")
		return
	}

	var req scaleNodePoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	provider, err := plugins.Get(cluster.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "provider not available")
		return
	}

	var creds map[string]string
	json.Unmarshal([]byte(cluster.Credentials), &creds)

	providerPoolID := pool.ProviderPoolID
	if providerPoolID == "" {
		providerPoolID = pool.Name
	}

	updated, err := provider.UpdateNodePool(r.Context(), creds, cluster.ProviderClusterID, providerPoolID, plugins.UpdateNodePoolOpts{
		NodeCount:   req.NodeCount,
		MinNodes:    req.MinNodes,
		MaxNodes:    req.MaxNodes,
		Autoscaling: req.Autoscaling,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to scale node pool: "+err.Error())
		return
	}

	// Update local record
	pool.NodeCount = updated.NodeCount
	pool.MinNodes = updated.MinNodes
	pool.MaxNodes = updated.MaxNodes
	pool.Autoscaling = updated.Autoscaling
	pool.Status = updated.Status
	h.clusters.UpdateNodePool(pool)

	msg := "Node pool '" + pool.Name + "' scaled"
	if req.NodeCount != nil {
		msg += " to " + strconv.Itoa(*req.NodeCount) + " nodes"
	}
	h.clusters.AddEvent(clusterID, "node_pool_scaled", msg, "{}")

	writeJSON(w, http.StatusOK, pool)
}

func (h *ClusterHandler) DeleteNodePool(w http.ResponseWriter, r *http.Request) {
	clusterID := chi.URLParam(r, "id")
	poolID := chi.URLParam(r, "poolId")

	cluster, err := h.clusters.GetByID(clusterID)
	if err != nil {
		writeError(w, http.StatusNotFound, "cluster not found")
		return
	}

	pool, err := h.clusters.GetNodePool(poolID)
	if err != nil {
		writeError(w, http.StatusNotFound, "node pool not found")
		return
	}

	provider, err := plugins.Get(cluster.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "provider not available")
		return
	}

	var creds map[string]string
	json.Unmarshal([]byte(cluster.Credentials), &creds)

	providerPoolID := pool.ProviderPoolID
	if providerPoolID == "" {
		providerPoolID = pool.Name
	}

	if err := provider.DeleteNodePool(r.Context(), creds, cluster.ProviderClusterID, providerPoolID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete node pool: "+err.Error())
		return
	}

	h.clusters.DeleteNodePool(poolID)
	h.clusters.AddEvent(clusterID, "node_pool_deleted", "Node pool '"+pool.Name+"' deleted", "{}")

	writeJSON(w, http.StatusOK, map[string]string{"message": "node pool deleted"})
}

// --- Nodes ---

func (h *ClusterHandler) ListNodes(w http.ResponseWriter, r *http.Request) {
	clusterID := chi.URLParam(r, "id")
	poolID := chi.URLParam(r, "poolId")

	cluster, err := h.clusters.GetByID(clusterID)
	if err != nil {
		writeError(w, http.StatusNotFound, "cluster not found")
		return
	}

	pool, err := h.clusters.GetNodePool(poolID)
	if err != nil {
		writeError(w, http.StatusNotFound, "node pool not found")
		return
	}

	provider, err := plugins.Get(cluster.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "provider not available")
		return
	}

	var creds map[string]string
	json.Unmarshal([]byte(cluster.Credentials), &creds)

	providerPoolID := pool.ProviderPoolID
	if providerPoolID == "" {
		providerPoolID = pool.Name
	}

	nodes, err := provider.ListNodes(r.Context(), creds, cluster.ProviderClusterID, providerPoolID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list nodes: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, nodes)
}

// --- Metrics ---

func (h *ClusterHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	clusterID := chi.URLParam(r, "id")
	cluster, err := h.clusters.GetByID(clusterID)
	if err != nil {
		writeError(w, http.StatusNotFound, "cluster not found")
		return
	}

	provider, err := plugins.Get(cluster.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "provider not available")
		return
	}

	var creds map[string]string
	json.Unmarshal([]byte(cluster.Credentials), &creds)

	metrics, err := provider.GetClusterMetrics(r.Context(), creds, cluster.ProviderClusterID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get metrics: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, metrics)
}

// --- Events ---

func (h *ClusterHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	clusterID := chi.URLParam(r, "id")
	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	events, err := h.clusters.ListEvents(clusterID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list events")
		return
	}
	writeJSON(w, http.StatusOK, events)
}

// --- Helpers ---

func (h *ClusterHandler) syncNodePools(ctx context.Context, cluster *models.Cluster, provider plugins.Provider, creds map[string]string) {
	pools, err := provider.ListNodePools(ctx, creds, cluster.ProviderClusterID)
	if err != nil {
		return
	}

	var dbPools []models.ClusterNodePool
	for _, p := range pools {
		labelsJSON, _ := json.Marshal(p.Labels)
		taintsJSON, _ := json.Marshal(p.Taints)
		dbPools = append(dbPools, models.ClusterNodePool{
			ClusterID:      cluster.ID,
			Name:           p.Name,
			ProviderPoolID: p.ID,
			MachineType:    p.MachineType,
			DiskSizeGB:     p.DiskSizeGB,
			DiskType:       p.DiskType,
			NodeCount:      p.NodeCount,
			MinNodes:       p.MinNodes,
			MaxNodes:       p.MaxNodes,
			Autoscaling:    p.Autoscaling,
			Preemptible:    p.Preemptible,
			SpotInstance:   p.SpotInstance,
			Status:         p.Status,
			K8sVersion:     p.K8sVersion,
			Labels:         string(labelsJSON),
			Taints:         string(taintsJSON),
		})
	}

	h.clusters.SyncNodePools(cluster.ID, dbPools)
}
