package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Cluster struct {
	ID                string     `db:"id" json:"id"`
	Name              string     `db:"name" json:"name"`
	DisplayName       string     `db:"display_name" json:"display_name"`
	Provider          string     `db:"provider" json:"provider"`
	ProviderClusterID string     `db:"provider_cluster_id" json:"provider_cluster_id"`
	ProjectID         string     `db:"project_id" json:"project_id"`
	Location          string     `db:"location" json:"location"`
	Status            string     `db:"status" json:"status"`
	K8sVersion        string     `db:"k8s_version" json:"k8s_version"`
	Endpoint          string     `db:"endpoint" json:"endpoint"`
	NodeCount         int        `db:"node_count" json:"node_count"`
	Credentials       string     `db:"credentials" json:"-"`
	Metadata          string     `db:"metadata" json:"metadata"`
	LastSyncedAt      *time.Time `db:"last_synced_at" json:"last_synced_at"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
}

type ClusterNodePool struct {
	ID             string    `db:"id" json:"id"`
	ClusterID      string    `db:"cluster_id" json:"cluster_id"`
	Name           string    `db:"name" json:"name"`
	ProviderPoolID string    `db:"provider_pool_id" json:"provider_pool_id"`
	MachineType    string    `db:"machine_type" json:"machine_type"`
	DiskSizeGB     int       `db:"disk_size_gb" json:"disk_size_gb"`
	DiskType       string    `db:"disk_type" json:"disk_type"`
	NodeCount      int       `db:"node_count" json:"node_count"`
	MinNodes       int       `db:"min_nodes" json:"min_nodes"`
	MaxNodes       int       `db:"max_nodes" json:"max_nodes"`
	Autoscaling    bool      `db:"autoscaling" json:"autoscaling"`
	Preemptible    bool      `db:"preemptible" json:"preemptible"`
	SpotInstance   bool      `db:"spot_instance" json:"spot_instance"`
	Status         string    `db:"status" json:"status"`
	K8sVersion     string    `db:"k8s_version" json:"k8s_version"`
	Labels         string    `db:"labels" json:"labels"`
	Taints         string    `db:"taints" json:"taints"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type ClusterEvent struct {
	ID        string    `db:"id" json:"id"`
	ClusterID string    `db:"cluster_id" json:"cluster_id"`
	EventType string    `db:"event_type" json:"event_type"`
	Message   string    `db:"message" json:"message"`
	Metadata  string    `db:"metadata" json:"metadata"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type ClusterStore struct {
	db *sqlx.DB
}

func NewClusterStore(db *sqlx.DB) *ClusterStore {
	return &ClusterStore{db: db}
}

// Cluster CRUD

func (s *ClusterStore) List() ([]Cluster, error) {
	var clusters []Cluster
	err := s.db.Select(&clusters, `SELECT * FROM clusters ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list clusters: %w", err)
	}
	return clusters, nil
}

func (s *ClusterStore) GetByID(id string) (*Cluster, error) {
	cluster := &Cluster{}
	err := s.db.Get(cluster, `SELECT * FROM clusters WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("get cluster: %w", err)
	}
	return cluster, nil
}

func (s *ClusterStore) Create(c *Cluster) (*Cluster, error) {
	result := &Cluster{}
	err := s.db.QueryRowx(
		`INSERT INTO clusters (name, display_name, provider, provider_cluster_id, project_id, location, status, k8s_version, endpoint, node_count, credentials, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *`,
		c.Name, c.DisplayName, c.Provider, c.ProviderClusterID, c.ProjectID, c.Location, c.Status, c.K8sVersion, c.Endpoint, c.NodeCount, c.Credentials, c.Metadata,
	).StructScan(result)
	if err != nil {
		return nil, fmt.Errorf("create cluster: %w", err)
	}
	return result, nil
}

func (s *ClusterStore) Update(c *Cluster) error {
	_, err := s.db.Exec(
		`UPDATE clusters SET name=?, display_name=?, status=?, k8s_version=?, endpoint=?, node_count=?, metadata=?, last_synced_at=CURRENT_TIMESTAMP, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		c.Name, c.DisplayName, c.Status, c.K8sVersion, c.Endpoint, c.NodeCount, c.Metadata, c.ID,
	)
	return err
}

func (s *ClusterStore) UpdateStatus(id, status string) error {
	_, err := s.db.Exec(`UPDATE clusters SET status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, status, id)
	return err
}

func (s *ClusterStore) UpdateCredentials(id, credentials string) error {
	_, err := s.db.Exec(`UPDATE clusters SET credentials=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, credentials, id)
	return err
}

func (s *ClusterStore) Delete(id string) error {
	_, err := s.db.Exec(`DELETE FROM clusters WHERE id = ?`, id)
	return err
}

func (s *ClusterStore) ListByProvider(provider string) ([]Cluster, error) {
	var clusters []Cluster
	err := s.db.Select(&clusters, `SELECT * FROM clusters WHERE provider = ? ORDER BY created_at DESC`, provider)
	if err != nil {
		return nil, fmt.Errorf("list clusters by provider: %w", err)
	}
	return clusters, nil
}

// Node Pool CRUD

func (s *ClusterStore) ListNodePools(clusterID string) ([]ClusterNodePool, error) {
	var pools []ClusterNodePool
	err := s.db.Select(&pools, `SELECT * FROM cluster_node_pools WHERE cluster_id = ? ORDER BY created_at`, clusterID)
	if err != nil {
		return nil, fmt.Errorf("list node pools: %w", err)
	}
	return pools, nil
}

func (s *ClusterStore) GetNodePool(id string) (*ClusterNodePool, error) {
	pool := &ClusterNodePool{}
	err := s.db.Get(pool, `SELECT * FROM cluster_node_pools WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("get node pool: %w", err)
	}
	return pool, nil
}

func (s *ClusterStore) CreateNodePool(p *ClusterNodePool) (*ClusterNodePool, error) {
	result := &ClusterNodePool{}
	err := s.db.QueryRowx(
		`INSERT INTO cluster_node_pools (cluster_id, name, provider_pool_id, machine_type, disk_size_gb, disk_type, node_count, min_nodes, max_nodes, autoscaling, preemptible, spot_instance, status, k8s_version, labels, taints)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *`,
		p.ClusterID, p.Name, p.ProviderPoolID, p.MachineType, p.DiskSizeGB, p.DiskType, p.NodeCount, p.MinNodes, p.MaxNodes, p.Autoscaling, p.Preemptible, p.SpotInstance, p.Status, p.K8sVersion, p.Labels, p.Taints,
	).StructScan(result)
	if err != nil {
		return nil, fmt.Errorf("create node pool: %w", err)
	}
	return result, nil
}

func (s *ClusterStore) UpdateNodePool(p *ClusterNodePool) error {
	_, err := s.db.Exec(
		`UPDATE cluster_node_pools SET node_count=?, min_nodes=?, max_nodes=?, autoscaling=?, status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		p.NodeCount, p.MinNodes, p.MaxNodes, p.Autoscaling, p.Status, p.ID,
	)
	return err
}

func (s *ClusterStore) DeleteNodePool(id string) error {
	_, err := s.db.Exec(`DELETE FROM cluster_node_pools WHERE id = ?`, id)
	return err
}

func (s *ClusterStore) SyncNodePools(clusterID string, pools []ClusterNodePool) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing pools
	tx.Exec(`DELETE FROM cluster_node_pools WHERE cluster_id = ?`, clusterID)

	// Insert fresh data
	for _, p := range pools {
		_, err := tx.Exec(
			`INSERT INTO cluster_node_pools (cluster_id, name, provider_pool_id, machine_type, disk_size_gb, disk_type, node_count, min_nodes, max_nodes, autoscaling, preemptible, spot_instance, status, k8s_version, labels, taints)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			clusterID, p.Name, p.ProviderPoolID, p.MachineType, p.DiskSizeGB, p.DiskType, p.NodeCount, p.MinNodes, p.MaxNodes, p.Autoscaling, p.Preemptible, p.SpotInstance, p.Status, p.K8sVersion, p.Labels, p.Taints,
		)
		if err != nil {
			return fmt.Errorf("sync node pool %s: %w", p.Name, err)
		}
	}

	return tx.Commit()
}

// Events

func (s *ClusterStore) AddEvent(clusterID, eventType, message, metadata string) error {
	_, err := s.db.Exec(
		`INSERT INTO cluster_events (cluster_id, event_type, message, metadata) VALUES (?, ?, ?, ?)`,
		clusterID, eventType, message, metadata,
	)
	return err
}

func (s *ClusterStore) ListEvents(clusterID string, limit int) ([]ClusterEvent, error) {
	var events []ClusterEvent
	err := s.db.Select(&events, `SELECT * FROM cluster_events WHERE cluster_id = ? ORDER BY created_at DESC LIMIT ?`, clusterID, limit)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	return events, nil
}
