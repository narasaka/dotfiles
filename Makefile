.PHONY: dev dev-api dev-web build docker-build helm-install helm-uninstall clean test

# Development
dev:
	@make -j2 dev-api dev-web

dev-api:
	KUBEPLOY_DEV=true go run ./cmd/kubeploy

dev-web:
	cd web && npm run dev

# Build
build:
	cd web && npm ci && npm run build
	CGO_ENABLED=0 go build -o bin/kubeploy ./cmd/kubeploy

# Test
test:
	go test ./...

# Docker
docker-build:
	docker build -t kubeploy:latest .

# Helm
helm-install:
	helm upgrade --install kubeploy ./deploy/helm/kubeploy \
		--namespace kubeploy-system --create-namespace

helm-uninstall:
	helm uninstall kubeploy --namespace kubeploy-system

helm-template:
	helm template kubeploy ./deploy/helm/kubeploy --namespace kubeploy-system

helm-lint:
	helm lint ./deploy/helm/kubeploy

# Clean
clean:
	rm -rf bin/ web/dist/ web/node_modules/
