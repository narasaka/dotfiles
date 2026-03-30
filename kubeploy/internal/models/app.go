package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type App struct {
	ID             string    `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	DisplayName    string    `db:"display_name" json:"display_name"`
	GitURL         string    `db:"git_url" json:"git_url"`
	GitBranch      string    `db:"git_branch" json:"git_branch"`
	GitSubpath     string    `db:"git_subpath" json:"git_subpath"`
	DockerfilePath string    `db:"dockerfile_path" json:"dockerfile_path"`
	RegistryImage  string    `db:"registry_image" json:"registry_image"`
	Namespace      string    `db:"namespace" json:"namespace"`
	Replicas       int       `db:"replicas" json:"replicas"`
	Port           int       `db:"port" json:"port"`
	EnvVars        string    `db:"env_vars" json:"env_vars"`
	AutoDeploy     bool      `db:"auto_deploy" json:"auto_deploy"`
	WebhookSecret  string    `db:"webhook_secret" json:"webhook_secret"`
	IngressHost    string    `db:"ingress_host" json:"ingress_host"`
	IngressTLS     bool      `db:"ingress_tls" json:"ingress_tls"`
	Status         string    `db:"status" json:"status"`
	CurrentBuildID *string   `db:"current_build_id" json:"current_build_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type AppStore struct {
	db *sqlx.DB
}

func NewAppStore(db *sqlx.DB) *AppStore {
	return &AppStore{db: db}
}

func (s *AppStore) List() ([]App, error) {
	var apps []App
	err := s.db.Select(&apps, `SELECT * FROM apps ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list apps: %w", err)
	}
	return apps, nil
}

func (s *AppStore) GetByID(id string) (*App, error) {
	app := &App{}
	err := s.db.Get(app, `SELECT * FROM apps WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("get app: %w", err)
	}
	return app, nil
}

func (s *AppStore) Create(app *App) (*App, error) {
	result := &App{}
	err := s.db.QueryRowx(
		`INSERT INTO apps (name, display_name, git_url, git_branch, git_subpath, dockerfile_path, registry_image, namespace, replicas, port, env_vars, auto_deploy, webhook_secret, ingress_host, ingress_tls)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *`,
		app.Name, app.DisplayName, app.GitURL, app.GitBranch, app.GitSubpath, app.DockerfilePath,
		app.RegistryImage, app.Namespace, app.Replicas, app.Port, app.EnvVars, app.AutoDeploy,
		app.WebhookSecret, app.IngressHost, app.IngressTLS,
	).StructScan(result)
	if err != nil {
		return nil, fmt.Errorf("create app: %w", err)
	}
	return result, nil
}

func (s *AppStore) Update(app *App) (*App, error) {
	result := &App{}
	err := s.db.QueryRowx(
		`UPDATE apps SET name=?, display_name=?, git_url=?, git_branch=?, git_subpath=?, dockerfile_path=?,
		registry_image=?, namespace=?, replicas=?, port=?, env_vars=?, auto_deploy=?, webhook_secret=?,
		ingress_host=?, ingress_tls=?, updated_at=CURRENT_TIMESTAMP WHERE id=? RETURNING *`,
		app.Name, app.DisplayName, app.GitURL, app.GitBranch, app.GitSubpath, app.DockerfilePath,
		app.RegistryImage, app.Namespace, app.Replicas, app.Port, app.EnvVars, app.AutoDeploy,
		app.WebhookSecret, app.IngressHost, app.IngressTLS, app.ID,
	).StructScan(result)
	if err != nil {
		return nil, fmt.Errorf("update app: %w", err)
	}
	return result, nil
}

func (s *AppStore) Delete(id string) error {
	_, err := s.db.Exec(`DELETE FROM apps WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete app: %w", err)
	}
	return nil
}

func (s *AppStore) UpdateStatus(id, status string) error {
	_, err := s.db.Exec(`UPDATE apps SET status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, status, id)
	return err
}

func (s *AppStore) UpdateCurrentBuild(id, buildID string) error {
	_, err := s.db.Exec(`UPDATE apps SET current_build_id=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, buildID, id)
	return err
}

func (s *AppStore) UpdateEnvVars(id, envVars string) error {
	_, err := s.db.Exec(`UPDATE apps SET env_vars=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`, envVars, id)
	return err
}

func (s *AppStore) GetByName(name string) (*App, error) {
	app := &App{}
	err := s.db.Get(app, `SELECT * FROM apps WHERE name = ?`, name)
	if err != nil {
		return nil, fmt.Errorf("get app by name: %w", err)
	}
	return app, nil
}
