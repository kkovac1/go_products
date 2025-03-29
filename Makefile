build:
	@go build -o bin/products cmd/main.go
run: build
	@./bin/products
test:
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@rm coverage.out