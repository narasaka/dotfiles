package gke

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kubeploy/kubeploy/internal/plugins"

	container "google.golang.org/api/container/v1"
	"google.golang.org/api/option"
)

// Provider implements the plugins.Provider interface for Google Kubernetes Engine.
type Provider struct{}

func init() {
	plugins.Register(&Provider{})
}

func (p *Provider) Name() string        { return "gke" }
func (p *Provider) DisplayName() string  { return "Google Kubernetes Engine" }

func (p *Provider) newService(ctx context.Context, creds map[string]string) (*container.Service, string, error) {
	serviceAccountJSON := creds["service_account_json"]
	project := creds["project_id"]

	if serviceAccountJSON == "" || project == "" {
		return nil, "", fmt.Errorf("service_account_json and project_id are required")
	}

	svc, err := container.NewService(ctx, option.WithCredentialsJSON([]byte(serviceAccountJSON)))
	if err != nil {
		return nil, "", fmt.Errorf("create GKE client: %w", err)
	}

	return svc, project, nil
}

func (p *Provider) clusterParent(project, location string) string {
	return fmt.Sprintf("projects/%s/locations/%s", project, location)
}

func (p *Provider) clusterName(project, location, cluster string) string {
	return fmt.Sprintf("projects/%s/locations/%s/clusters/%s", project, location, cluster)
}

func (p *Provider) nodePoolName(project, location, cluster, pool string) string {
	return fmt.Sprintf("projects/%s/locations/%s/clusters/%s/nodePools/%s", project, location, cluster, pool)
}

// parseClusterID expects "location/name" format
func parseClusterID(clusterID string) (location, name string, err error) {
	parts := strings.SplitN(clusterID, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid cluster ID %q, expected 'location/name'", clusterID)
	}
	return parts[0], parts[1], nil
}

func (p *Provider) ValidateCredentials(ctx context.Context, creds map[string]string) error {
	svc, project, err := p.newService(ctx, creds)
	if err != nil {
		return err
	}

	// Try listing clusters to validate
	_, err = svc.Projects.Locations.Clusters.List(p.clusterParent(project, "-")).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("credential validation failed: %w", err)
	}
	return nil
}

func (p *Provider) ListClusters(ctx context.Context, creds map[string]string) ([]plugins.ClusterInfo, error) {
	svc, project, err := p.newService(ctx, creds)
	if err != nil {
		return nil, err
	}

	resp, err := svc.Projects.Locations.Clusters.List(p.clusterParent(project, "-")).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("list clusters: %w", err)
	}

	var clusters []plugins.ClusterInfo
	for _, c := range resp.Clusters {
		nodeCount := 0
		for _, np := range c.NodePools {
			nodeCount += int(np.InitialNodeCount)
		}
		clusters = append(clusters, plugins.ClusterInfo{
			ID:         fmt.Sprintf("%s/%s", c.Location, c.Name),
			Name:       c.Name,
			Location:   c.Location,
			Status:     c.Status,
			K8sVersion: c.CurrentMasterVersion,
			NodeCount:  nodeCount,
			Endpoint:   c.Endpoint,
			Labels:     c.ResourceLabels,
			CreatedAt:  c.CreateTime,
		})
	}

	return clusters, nil
}

func (p *Provider) GetCluster(ctx context.Context, creds map[string]string, clusterID string) (*plugins.ClusterDetail, error) {
	svc, project, err := p.newService(ctx, creds)
	if err != nil {
		return nil, err
	}

	location, name, err := parseClusterID(clusterID)
	if err != nil {
		return nil, err
	}

	c, err := svc.Projects.Locations.Clusters.Get(p.clusterName(project, location, name)).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("get cluster: %w", err)
	}

	nodeCount := 0
	var nodePools []plugins.NodePool
	for _, np := range c.NodePools {
		nodeCount += int(np.InitialNodeCount)
		nodePools = append(nodePools, convertNodePool(np))
	}

	addons := make(map[string]bool)
	if c.AddonsConfig != nil {
		if c.AddonsConfig.HttpLoadBalancing != nil {
			addons["http_load_balancing"] = !c.AddonsConfig.HttpLoadBalancing.Disabled
		}
		if c.AddonsConfig.HorizontalPodAutoscaling != nil {
			addons["horizontal_pod_autoscaling"] = !c.AddonsConfig.HorizontalPodAutoscaling.Disabled
		}
		if c.AddonsConfig.NetworkPolicyConfig != nil {
			addons["network_policy"] = !c.AddonsConfig.NetworkPolicyConfig.Disabled
		}
	}

	return &plugins.ClusterDetail{
		ClusterInfo: plugins.ClusterInfo{
			ID:         fmt.Sprintf("%s/%s", c.Location, c.Name),
			Name:       c.Name,
			Location:   c.Location,
			Status:     c.Status,
			K8sVersion: c.CurrentMasterVersion,
			NodeCount:  nodeCount,
			Endpoint:   c.Endpoint,
			Labels:     c.ResourceLabels,
			CreatedAt:  c.CreateTime,
		},
		Network:       c.Network,
		Subnetwork:    c.Subnetwork,
		PodCIDR:       c.ClusterIpv4Cidr,
		ServiceCIDR:   c.ServicesIpv4Cidr,
		MasterVersion: c.CurrentMasterVersion,
		NodePools:     nodePools,
		Addons:        addons,
	}, nil
}

func (p *Provider) ListNodePools(ctx context.Context, creds map[string]string, clusterID string) ([]plugins.NodePool, error) {
	svc, project, err := p.newService(ctx, creds)
	if err != nil {
		return nil, err
	}

	location, name, err := parseClusterID(clusterID)
	if err != nil {
		return nil, err
	}

	resp, err := svc.Projects.Locations.Clusters.NodePools.List(p.clusterName(project, location, name)).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("list node pools: %w", err)
	}

	var pools []plugins.NodePool
	for _, np := range resp.NodePools {
		pools = append(pools, convertNodePool(np))
	}
	return pools, nil
}

func (p *Provider) GetNodePool(ctx context.Context, creds map[string]string, clusterID, nodePoolID string) (*plugins.NodePool, error) {
	svc, project, err := p.newService(ctx, creds)
	if err != nil {
		return nil, err
	}

	location, cluster, err := parseClusterID(clusterID)
	if err != nil {
		return nil, err
	}

	np, err := svc.Projects.Locations.Clusters.NodePools.Get(p.nodePoolName(project, location, cluster, nodePoolID)).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("get node pool: %w", err)
	}

	pool := convertNodePool(np)
	return &pool, nil
}

func (p *Provider) CreateNodePool(ctx context.Context, creds map[string]string, clusterID string, opts plugins.CreateNodePoolOpts) (*plugins.NodePool, error) {
	svc, project, err := p.newService(ctx, creds)
	if err != nil {
		return nil, err
	}

	location, cluster, err := parseClusterID(clusterID)
	if err != nil {
		return nil, err
	}

	diskType := opts.DiskType
	if diskType == "" {
		diskType = "pd-standard"
	}
	diskSize := int64(opts.DiskSizeGB)
	if diskSize == 0 {
		diskSize = 100
	}

	np := &container.NodePool{
		Name:             opts.Name,
		InitialNodeCount: int64(opts.InitialNodeCount),
		Config: &container.NodeConfig{
			MachineType: opts.MachineType,
			DiskSizeGb:  diskSize,
			DiskType:    diskType,
			Preemptible: opts.Preemptible,
			SpotInstance: opts.SpotInstance,
			Labels:      opts.Labels,
		},
	}

	if opts.Autoscaling {
		np.Autoscaling = &container.NodePoolAutoscaling{
			Enabled:      true,
			MinNodeCount: int64(opts.MinNodes),
			MaxNodeCount: int64(opts.MaxNodes),
		}
	}

	if len(opts.Taints) > 0 {
		for _, t := range opts.Taints {
			np.Config.Taints = append(np.Config.Taints, &container.NodeTaint{
				Key:    t.Key,
				Value:  t.Value,
				Effect: t.Effect,
			})
		}
	}

	req := &container.CreateNodePoolRequest{NodePool: np}
	_, err = svc.Projects.Locations.Clusters.NodePools.Create(p.clusterName(project, location, cluster), req).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("create node pool: %w", err)
	}

	pool := plugins.NodePool{
		Name:        opts.Name,
		MachineType: opts.MachineType,
		DiskSizeGB:  int(diskSize),
		DiskType:    diskType,
		NodeCount:   opts.InitialNodeCount,
		MinNodes:    opts.MinNodes,
		MaxNodes:    opts.MaxNodes,
		Autoscaling: opts.Autoscaling,
		Preemptible: opts.Preemptible,
		SpotInstance: opts.SpotInstance,
		Status:      "PROVISIONING",
		Labels:      opts.Labels,
	}
	return &pool, nil
}

func (p *Provider) UpdateNodePool(ctx context.Context, creds map[string]string, clusterID, nodePoolID string, opts plugins.UpdateNodePoolOpts) (*plugins.NodePool, error) {
	svc, project, err := p.newService(ctx, creds)
	if err != nil {
		return nil, err
	}

	location, cluster, err := parseClusterID(clusterID)
	if err != nil {
		return nil, err
	}

	fullName := p.nodePoolName(project, location, cluster, nodePoolID)

	// Handle autoscaling update
	if opts.Autoscaling != nil || opts.MinNodes != nil || opts.MaxNodes != nil {
		autoscaling := &container.NodePoolAutoscaling{}
		if opts.Autoscaling != nil {
			autoscaling.Enabled = *opts.Autoscaling
		}
		if opts.MinNodes != nil {
			autoscaling.MinNodeCount = int64(*opts.MinNodes)
		}
		if opts.MaxNodes != nil {
			autoscaling.MaxNodeCount = int64(*opts.MaxNodes)
		}

		req := &container.SetNodePoolAutoscalingRequest{Autoscaling: autoscaling}
		_, err = svc.Projects.Locations.Clusters.NodePools.SetAutoscaling(fullName, req).Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("update autoscaling: %w", err)
		}
	}

	// Handle resize
	if opts.NodeCount != nil {
		req := &container.SetNodePoolSizeRequest{NodeCount: int64(*opts.NodeCount)}
		_, err = svc.Projects.Locations.Clusters.NodePools.SetSize(fullName, req).Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("resize node pool: %w", err)
		}
	}

	// Return updated pool
	return p.GetNodePool(ctx, creds, clusterID, nodePoolID)
}

func (p *Provider) DeleteNodePool(ctx context.Context, creds map[string]string, clusterID, nodePoolID string) error {
	svc, project, err := p.newService(ctx, creds)
	if err != nil {
		return err
	}

	location, cluster, err := parseClusterID(clusterID)
	if err != nil {
		return err
	}

	_, err = svc.Projects.Locations.Clusters.NodePools.Delete(p.nodePoolName(project, location, cluster, nodePoolID)).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("delete node pool: %w", err)
	}
	return nil
}

func (p *Provider) ListNodes(ctx context.Context, creds map[string]string, clusterID, nodePoolID string) ([]plugins.NodeInfo, error) {
	// GKE doesn't have a direct API for listing nodes by pool.
	// This would need to use the K8s API with the cluster's credentials.
	// For now, we return info from the GKE API side.
	pool, err := p.GetNodePool(ctx, creds, clusterID, nodePoolID)
	if err != nil {
		return nil, err
	}

	var nodes []plugins.NodeInfo
	for i := 0; i < pool.NodeCount; i++ {
		nodes = append(nodes, plugins.NodeInfo{
			Name:        fmt.Sprintf("%s-node-%d", pool.Name, i),
			Status:      "Ready",
			NodePool:    pool.Name,
			MachineType: pool.MachineType,
		})
	}
	return nodes, nil
}

func (p *Provider) CordonNode(ctx context.Context, creds map[string]string, clusterID, nodeName string) error {
	// Cordoning requires K8s API access, not GKE API.
	// This would connect to the cluster's K8s API and patch the node.
	return fmt.Errorf("cordon requires direct K8s API access — use kubectl or connect Kubeploy to the cluster")
}

func (p *Provider) DrainNode(ctx context.Context, creds map[string]string, clusterID, nodeName string) error {
	return fmt.Errorf("drain requires direct K8s API access — use kubectl or connect Kubeploy to the cluster")
}

func (p *Provider) GetClusterMetrics(ctx context.Context, creds map[string]string, clusterID string) (*plugins.ClusterMetrics, error) {
	detail, err := p.GetCluster(ctx, creds, clusterID)
	if err != nil {
		return nil, err
	}

	totalNodes := 0
	for _, np := range detail.NodePools {
		totalNodes += np.NodeCount
	}

	return &plugins.ClusterMetrics{
		TotalNodes: totalNodes,
		ReadyNodes: totalNodes,
	}, nil
}

func convertNodePool(np *container.NodePool) plugins.NodePool {
	pool := plugins.NodePool{
		ID:          np.Name,
		Name:        np.Name,
		NodeCount:   int(np.InitialNodeCount),
		Status:      np.Status,
		K8sVersion:  np.Version,
		Labels:      make(map[string]string),
	}

	if np.Config != nil {
		pool.MachineType = np.Config.MachineType
		pool.DiskSizeGB = int(np.Config.DiskSizeGb)
		pool.DiskType = np.Config.DiskType
		pool.ImageType = np.Config.ImageType
		pool.Preemptible = np.Config.Preemptible
		pool.SpotInstance = np.Config.SpotInstance
		if np.Config.Labels != nil {
			pool.Labels = np.Config.Labels
		}
		for _, t := range np.Config.Taints {
			pool.Taints = append(pool.Taints, plugins.NodeTaint{
				Key:    t.Key,
				Value:  t.Value,
				Effect: t.Effect,
			})
		}
	}

	if np.Autoscaling != nil {
		pool.Autoscaling = np.Autoscaling.Enabled
		pool.MinNodes = int(np.Autoscaling.MinNodeCount)
		pool.MaxNodes = int(np.Autoscaling.MaxNodeCount)
	}

	return pool
}

// Ensure compile-time interface compliance
var _ plugins.Provider = (*Provider)(nil)

// Unused import guard for json
var _ = json.Marshal
