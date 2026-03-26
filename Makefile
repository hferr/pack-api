# Build and run app in Docker
build-run:
	docker-compose up --build

# Tidy
tidy:
	go mod tidy
