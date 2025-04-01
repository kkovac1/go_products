FROM golang:1.24.1 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the API binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

# Build the migration binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate ./cmd/migrate/main.go

FROM alpine:latest AS build-release-stage
WORKDIR /
COPY --from=build-stage /api /api
COPY --from=build-stage /migrate /migrate
COPY cmd/migrate/migrations /cmd/migrate/migrations
RUN chmod +x /api  # Ensure executable permissions
EXPOSE 8080
CMD ["/api"]