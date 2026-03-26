# Build and run app in Docker
build-run:
	docker-compose up --build

# Run app
run:
	docker-compose up -d

# Run tests locally
run-tests:
	go test ./... -race

# Tidy
tidy:
	go mod tidy
