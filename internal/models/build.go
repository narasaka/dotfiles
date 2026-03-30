package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Build struct {
	ID            string     `db:"id" json:"id"`
	AppID         string     `db:"app_id" json:"app_id"`
	CommitSHA     string     `db:"commit_sha" json:"commit_sha"`
	CommitMessage string     `db:"commit_message" json:"commit_message"`
	CommitAuthor  string     `db:"commit_author" json:"commit_author"`
	ImageTag      string     `db:"image_tag" json:"image_tag"`
	Status        string     `db:"status" json:"status"`
	BuildJobName  string     `db:"build_job_name" json:"build_job_name"`
	Logs          string     `db:"logs" json:"logs"`
	StartedAt     *time.Time `db:"started_at" json:"started_at"`
	FinishedAt    *time.Time `db:"finished_at" json:"finished_at"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}

type BuildStore struct {
	db *sqlx.DB
}

func NewBuildStore(db *sqlx.DB) *BuildStore {
	return &BuildStore{db: db}
}

func (s *BuildStore) ListByApp(appID string) ([]Build, error) {
	var builds []Build
	err := s.db.Select(&builds, `SELECT * FROM builds WHERE app_id = ? ORDER BY created_at DESC`, appID)
	if err != nil {
		return nil, fmt.Errorf("list builds: %w", err)
	}
	return builds, nil
}

func (s *BuildStore) GetByID(id string) (*Build, error) {
	build := &Build{}
	err := s.db.Get(build, `SELECT * FROM builds WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("get build: %w", err)
	}
	return build, nil
}

func (s *BuildStore) Create(build *Build) (*Build, error) {
	result := &Build{}
	err := s.db.QueryRowx(
		`INSERT INTO builds (app_id, commit_sha, commit_message, commit_author, image_tag, status, build_job_name)
		VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING *`,
		build.AppID, build.CommitSHA, build.CommitMessage, build.CommitAuthor,
		build.ImageTag, build.Status, build.BuildJobName,
	).StructScan(result)
	if err != nil {
		return nil, fmt.Errorf("create build: %w", err)
	}
	return result, nil
}

func (s *BuildStore) UpdateStatus(id, status string) error {
	_, err := s.db.Exec(`UPDATE builds SET status=? WHERE id=?`, status, id)
	return err
}

func (s *BuildStore) SetStarted(id string) error {
	_, err := s.db.Exec(`UPDATE builds SET status='building', started_at=CURRENT_TIMESTAMP WHERE id=?`, id)
	return err
}

func (s *BuildStore) SetFinished(id, status string) error {
	_, err := s.db.Exec(`UPDATE builds SET status=?, finished_at=CURRENT_TIMESTAMP WHERE id=?`, status, id)
	return err
}

func (s *BuildStore) AppendLogs(id, logLine string) error {
	_, err := s.db.Exec(`UPDATE builds SET logs = logs || ? WHERE id=?`, logLine, id)
	return err
}

func (s *BuildStore) SetBuildJobName(id, jobName string) error {
	_, err := s.db.Exec(`UPDATE builds SET build_job_name=? WHERE id=?`, jobName, id)
	return err
}
