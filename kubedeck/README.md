# Kubedeck

A self-hosted PaaS for Kubernetes. Deploy from git push with a clean web UI.

## Features

- **Git push to deploy** - Automatic builds and deployments via GitHub/GitLab webhooks
- **Web UI** - Manage apps, builds, and deployments from a dark, terminal-inspired dashboard
- **Real-time logs** - Stream build and runtime logs via WebSocket (xterm.js)
- **Rollback** - One-click rollback to any previous deployment
- **In-cluster builds** - Kaniko-based container builds, no Docker daemon needed
- **Zero external deps** - SQLite database, single binary, single Helm install
- **Ingress management** - Automatic Service and Ingress creation with optional TLS

## Architecture

```
Web UI (React) <-> API Server (Go/Chi) <-> K8s API
                                        |
                        +---------------+---------------+
                        |               |               |
                  Build Controller  Deploy Controller  Log Streamer
                    (Kaniko)        (K8s Resources)    (WebSocket)
                        |
                    SQLite DB
```

## Quickstart

### Helm Install

```bash
helm upgrade --install kubedeck ./deploy/helm/kubedeck \
  --namespace kubedeck-system \
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
| `KUBEDECK_DEV` | `false` | Enable dev mode |
| `KUBEDECK_PORT` | `8080` | API server port |
| `KUBEDECK_DB_PATH` | `./kubedeck.db` | SQLite database path |
| `KUBEDECK_SESSION_SECRET` | auto-generated | Session encryption key |
| `KUBEDECK_NAMESPACE` | `kubedeck-system` | Kubedeck's namespace |
| `KUBEDECK_APP_NAMESPACE` | `kubedeck-apps` | Default app namespace |
| `KUBEDECK_REGISTRY_URL` | | Container registry URL |
| `KUBEDECK_REGISTRY_USERNAME` | | Registry username |
| `KUBEDECK_REGISTRY_PASSWORD` | | Registry password |

## Tech Stack

**Backend:** Go 1.22+, Chi router, SQLite (modernc.org/sqlite), client-go, nhooyr.io/websocket

**Frontend:** React 18, TypeScript, Vite, Tailwind CSS, Zustand, xterm.js, Lucide icons

**Infrastructure:** Kaniko (in-cluster builds), Helm chart, standard K8s resources (no CRDs)

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
