package plugins

import "context"

// Provider is the interface that all cluster provider plugins must implement.
// Currently supports GKE, with EKS, AKS, and bare-metal planned for future.
type Provider interface {
	// Name returns the provider identifier (e.g. "gke", "eks", "aks")
	Name() string

	// DisplayName returns a human-readable name (e.g. "Google Kubernetes Engine")
	DisplayName() string

	// ValidateCredentials checks if the provided credentials are valid
	ValidateCredentials(ctx context.Context, creds map[string]string) error

	// ListClusters returns all clusters accessible with the configured credentials
	ListClusters(ctx context.Context, creds map[string]string) ([]ClusterInfo, error)

	// GetCluster returns detailed info about a specific cluster
	GetCluster(ctx context.Context, creds map[string]string, clusterID string) (*ClusterDetail, error)

	// ListNodePools returns all node pools for a cluster
	ListNodePools(ctx context.Context, creds map[string]string, clusterID string) ([]NodePool, error)

	// GetNodePool returns details of a specific node pool
	GetNodePool(ctx context.Context, creds map[string]string, clusterID, nodePoolID string) (*NodePool, error)

	// CreateNodePool provisions a new node pool in the cluster
	CreateNodePool(ctx context.Context, creds map[string]string, clusterID string, opts CreateNodePoolOpts) (*NodePool, error)

	// UpdateNodePool modifies an existing node pool (resize, autoscaling, etc.)
	UpdateNodePool(ctx context.Context, creds map[string]string, clusterID, nodePoolID string, opts UpdateNodePoolOpts) (*NodePool, error)

	// DeleteNodePool removes a node pool from the cluster
	DeleteNodePool(ctx context.Context, creds map[string]string, clusterID, nodePoolID string) error

	// ListNodes returns all nodes in a specific node pool
	ListNodes(ctx context.Context, creds map[string]string, clusterID, nodePoolID string) ([]NodeInfo, error)

	// CordonNode marks a node as unschedulable
	CordonNode(ctx context.Context, creds map[string]string, clusterID, nodeName string) error

	// DrainNode safely evicts pods from a node
	DrainNode(ctx context.Context, creds map[string]string, clusterID, nodeName string) error

	// GetClusterMetrics returns resource usage metrics for the cluster
	GetClusterMetrics(ctx context.Context, creds map[string]string, clusterID string) (*ClusterMetrics, error)
}

type ClusterInfo struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Location        string            `json:"location"`
	Status          string            `json:"status"`
	K8sVersion      string            `json:"k8s_version"`
	NodeCount       int               `json:"node_count"`
	Endpoint        string            `json:"endpoint"`
	Labels          map[string]string `json:"labels"`
	CreatedAt       string            `json:"created_at"`
}

type ClusterDetail struct {
	ClusterInfo
	Network         string            `json:"network"`
	Subnetwork      string            `json:"subnetwork"`
	PodCIDR         string            `json:"pod_cidr"`
	ServiceCIDR     string            `json:"service_cidr"`
	MasterVersion   string            `json:"master_version"`
	NodePools       []NodePool        `json:"node_pools"`
	Addons          map[string]bool   `json:"addons"`
}

type NodePool struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	MachineType     string            `json:"machine_type"`
	DiskSizeGB      int               `json:"disk_size_gb"`
	DiskType        string            `json:"disk_type"`
	ImageType       string            `json:"image_type"`
	NodeCount       int               `json:"node_count"`
	MinNodes        int               `json:"min_nodes"`
	MaxNodes        int               `json:"max_nodes"`
	Autoscaling     bool              `json:"autoscaling"`
	Preemptible     bool              `json:"preemptible"`
	SpotInstance    bool              `json:"spot_instance"`
	Status          string            `json:"status"`
	K8sVersion      string            `json:"k8s_version"`
	Labels          map[string]string `json:"labels"`
	Taints          []NodeTaint       `json:"taints"`
}

type NodeTaint struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"`
}

type CreateNodePoolOpts struct {
	Name            string            `json:"name"`
	MachineType     string            `json:"machine_type"`
	DiskSizeGB      int               `json:"disk_size_gb"`
	DiskType        string            `json:"disk_type"`
	InitialNodeCount int              `json:"initial_node_count"`
	MinNodes        int               `json:"min_nodes"`
	MaxNodes        int               `json:"max_nodes"`
	Autoscaling     bool              `json:"autoscaling"`
	Preemptible     bool              `json:"preemptible"`
	SpotInstance    bool              `json:"spot_instance"`
	Labels          map[string]string `json:"labels"`
	Taints          []NodeTaint       `json:"taints"`
}

type UpdateNodePoolOpts struct {
	NodeCount       *int              `json:"node_count,omitempty"`
	MinNodes        *int              `json:"min_nodes,omitempty"`
	MaxNodes        *int              `json:"max_nodes,omitempty"`
	Autoscaling     *bool             `json:"autoscaling,omitempty"`
}

type NodeInfo struct {
	Name            string            `json:"name"`
	Status          string            `json:"status"`
	NodePool        string            `json:"node_pool"`
	MachineType     string            `json:"machine_type"`
	Zone            string            `json:"zone"`
	K8sVersion      string            `json:"k8s_version"`
	CPUCapacity     string            `json:"cpu_capacity"`
	MemoryCapacity  string            `json:"memory_capacity"`
	CPUAllocatable  string            `json:"cpu_allocatable"`
	MemoryAllocatable string          `json:"memory_allocatable"`
	PodCount        int               `json:"pod_count"`
	Conditions      []NodeCondition   `json:"conditions"`
	CreatedAt       string            `json:"created_at"`
}

type NodeCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ClusterMetrics struct {
	TotalCPU        string `json:"total_cpu"`
	UsedCPU         string `json:"used_cpu"`
	TotalMemory     string `json:"total_memory"`
	UsedMemory      string `json:"used_memory"`
	TotalPods       int    `json:"total_pods"`
	RunningPods     int    `json:"running_pods"`
	TotalNodes      int    `json:"total_nodes"`
	ReadyNodes      int    `json:"ready_nodes"`
}
