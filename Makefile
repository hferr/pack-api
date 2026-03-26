# Build and run app in Docker
build-run:
	docker-compose up --build

# Run app
run:
	docker-compose up -d

# Tidy
tidy:
	go mod tidy
