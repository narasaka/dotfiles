# Kubeploy

A self-hosted PaaS for Kubernetes. Deploy from git push with a clean web UI.

## Features

- **Git push to deploy** - Automatic builds and deployments via GitHub/GitLab webhooks
- **Web UI** - Manage apps, builds, and deployments from a dark, terminal-inspired dashboard
- **Real-time logs** - Stream build and runtime logs via WebSocket (xterm.js)
- **Rollback** - One-click rollback to any previous deployment
- **In-cluster builds** - BuildKit-powered container builds, no Docker daemon needed
- **Zero external deps** - SQLite database, single binary, single Helm install
- **Ingress management** - Automatic Service and Ingress creation with optional TLS

## Architecture

```
Web UI (React) <-> API Server (Go/Chi) <-> K8s API
                                        |
                        +---------------+---------------+
                        |               |               |
                  Build Controller  Deploy Controller  Log Streamer
                   (BuildKit)       (K8s Resources)    (WebSocket)
                        |
                    SQLite DB
```

## Quickstart

### Helm Install

```bash
helm upgrade --install kubeploy ./deploy/helm/kubeploy \
  --namespace kubeploy-system \
  --create-namespace \
  --set registry.url=ghcr.io \
  --set registry.username=YOUR_USER \
  --set registry.password=YOUR_TOKEN
```

### Local Development

```bash
# Start both API server and frontend dev server
make dev

# API runs on :8080, frontend on :5173 (proxies API)
```

### Docker Build

```bash
make docker-build
```

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `KUBEPLOY_DEV` | `false` | Enable dev mode |
| `KUBEPLOY_PORT` | `8080` | API server port |
| `KUBEPLOY_DB_PATH` | `./kubeploy.db` | SQLite database path |
| `KUBEPLOY_SESSION_SECRET` | auto-generated | Session encryption key |
| `KUBEPLOY_NAMESPACE` | `kubeploy-system` | Kubeploy's namespace |
| `KUBEPLOY_APP_NAMESPACE` | `kubeploy-apps` | Default app namespace |
| `KUBEPLOY_REGISTRY_URL` | | Container registry URL |
| `KUBEPLOY_REGISTRY_USERNAME` | | Registry username |
| `KUBEPLOY_REGISTRY_PASSWORD` | | Registry password |
| `KUBEPLOY_BUILDKIT_ADDR` | `tcp://kubeploy-buildkitd:1234` | BuildKit daemon address |

## Tech Stack

**Backend:** Go 1.22+, Chi router, SQLite (modernc.org/sqlite), client-go, nhooyr.io/websocket

**Frontend:** React 18, TypeScript, Vite, Tailwind CSS, Zustand, xterm.js, Lucide icons

**Infrastructure:** BuildKit (in-cluster builds), Helm chart, standard K8s resources (no CRDs)

## API

All routes under `/api/v1`. Auth via session cookie.

- `POST /auth/setup` - Create initial admin user
- `POST /auth/login` - Login
- `GET /apps` - List apps
- `POST /apps` - Create app
- `POST /apps/:id/builds` - Trigger build
- `WS /apps/:id/logs/ws` - Stream runtime logs
- `WS /builds/:id/logs/ws` - Stream build logs
- `POST /webhooks/github/:app_id` - GitHub push webhook
- `POST /webhooks/gitlab/:app_id` - GitLab push webhook

See full API reference in the source code.

## License

Apache 2.0
