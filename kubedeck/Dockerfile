# Stage 1: Build frontend
FROM node:20-alpine AS frontend
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ .
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.22-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o /kubedeck ./cmd/kubedeck

# Stage 3: Final image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates git
COPY --from=backend /kubedeck /usr/local/bin/kubedeck
EXPOSE 8080
ENTRYPOINT ["kubedeck"]
