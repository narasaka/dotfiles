.PHONY: dev dev-api dev-web build docker-build helm-install helm-uninstall clean test

# Development
dev:
	@make -j2 dev-api dev-web

dev-api:
	KUBEDECK_DEV=true go run ./cmd/kubedeck

dev-web:
	cd web && npm run dev

# Build
build:
	cd web && npm ci && npm run build
	CGO_ENABLED=0 go build -o bin/kubedeck ./cmd/kubedeck

# Test
test:
	go test ./...

# Docker
docker-build:
	docker build -t kubedeck:latest .

# Helm
helm-install:
	helm upgrade --install kubedeck ./deploy/helm/kubedeck \
		--namespace kubedeck-system --create-namespace

helm-uninstall:
	helm uninstall kubedeck --namespace kubedeck-system

helm-template:
	helm template kubedeck ./deploy/helm/kubedeck --namespace kubedeck-system

helm-lint:
	helm lint ./deploy/helm/kubedeck

# Clean
clean:
	rm -rf bin/ web/dist/ web/node_modules/
