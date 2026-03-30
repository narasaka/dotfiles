package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Deployment struct {
	ID                string    `db:"id" json:"id"`
	AppID             string    `db:"app_id" json:"app_id"`
	BuildID           string    `db:"build_id" json:"build_id"`
	K8sDeploymentName string    `db:"k8s_deployment_name" json:"k8s_deployment_name"`
	ReplicasDesired   int       `db:"replicas_desired" json:"replicas_desired"`
	ReplicasReady     int       `db:"replicas_ready" json:"replicas_ready"`
	Status            string    `db:"status" json:"status"`
	RolledBackTo      *string   `db:"rolled_back_to" json:"rolled_back_to"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
}

type DeploymentStore struct {
	db *sqlx.DB
}

func NewDeploymentStore(db *sqlx.DB) *DeploymentStore {
	return &DeploymentStore{db: db}
}

func (s *DeploymentStore) ListByApp(appID string) ([]Deployment, error) {
	var deployments []Deployment
	err := s.db.Select(&deployments, `SELECT * FROM deployments WHERE app_id = ? ORDER BY created_at DESC`, appID)
	if err != nil {
		return nil, fmt.Errorf("list deployments: %w", err)
	}
	return deployments, nil
}

func (s *DeploymentStore) GetByID(id string) (*Deployment, error) {
	dep := &Deployment{}
	err := s.db.Get(dep, `SELECT * FROM deployments WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("get deployment: %w", err)
	}
	return dep, nil
}

func (s *DeploymentStore) Create(dep *Deployment) (*Deployment, error) {
	result := &Deployment{}
	err := s.db.QueryRowx(
		`INSERT INTO deployments (app_id, build_id, k8s_deployment_name, replicas_desired, replicas_ready, status, rolled_back_to)
		VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING *`,
		dep.AppID, dep.BuildID, dep.K8sDeploymentName, dep.ReplicasDesired, dep.ReplicasReady, dep.Status, dep.RolledBackTo,
	).StructScan(result)
	if err != nil {
		return nil, fmt.Errorf("create deployment: %w", err)
	}
	return result, nil
}

func (s *DeploymentStore) UpdateStatus(id, status string) error {
	_, err := s.db.Exec(`UPDATE deployments SET status=? WHERE id=?`, status, id)
	return err
}

func (s *DeploymentStore) UpdateReplicas(id string, ready int) error {
	_, err := s.db.Exec(`UPDATE deployments SET replicas_ready=? WHERE id=?`, ready, id)
	return err
}

func (s *DeploymentStore) GetLatestByApp(appID string) (*Deployment, error) {
	dep := &Deployment{}
	err := s.db.Get(dep, `SELECT * FROM deployments WHERE app_id = ? ORDER BY created_at DESC LIMIT 1`, appID)
	if err != nil {
		return nil, fmt.Errorf("get latest deployment: %w", err)
	}
	return dep, nil
}
