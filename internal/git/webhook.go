package git

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type WebhookPayload struct {
	Branch        string
	CommitSHA     string
	CommitMessage string
	CommitAuthor  string
}

func ParseGitHubWebhook(r *http.Request, secret string) (*WebhookPayload, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	defer r.Body.Close()

	// Validate signature if secret is set
	if secret != "" {
		sig := r.Header.Get("X-Hub-Signature-256")
		if sig == "" {
			return nil, fmt.Errorf("missing signature header")
		}
		sig = strings.TrimPrefix(sig, "sha256=")
		if !validateHMAC(body, secret, sig) {
			return nil, fmt.Errorf("invalid signature")
		}
	}

	var payload struct {
		Ref        string `json:"ref"`
		HeadCommit struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"head_commit"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("parse payload: %w", err)
	}

	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

	return &WebhookPayload{
		Branch:        branch,
		CommitSHA:     payload.HeadCommit.ID,
		CommitMessage: payload.HeadCommit.Message,
		CommitAuthor:  payload.HeadCommit.Author.Name,
	}, nil
}

func ParseGitLabWebhook(r *http.Request, secret string) (*WebhookPayload, error) {
	// Validate token
	if secret != "" {
		token := r.Header.Get("X-Gitlab-Token")
		if token != secret {
			return nil, fmt.Errorf("invalid token")
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	defer r.Body.Close()

	var payload struct {
		Ref     string `json:"ref"`
		Commits []struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"commits"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("parse payload: %w", err)
	}

	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

	var commitSHA, commitMessage, commitAuthor string
	if len(payload.Commits) > 0 {
		last := payload.Commits[len(payload.Commits)-1]
		commitSHA = last.ID
		commitMessage = last.Message
		commitAuthor = last.Author.Name
	}

	return &WebhookPayload{
		Branch:        branch,
		CommitSHA:     commitSHA,
		CommitMessage: commitMessage,
		CommitAuthor:  commitAuthor,
	}, nil
}

func validateHMAC(body []byte, secret, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
