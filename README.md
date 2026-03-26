# Pack API

Golang API for managing Packs, using PostgreSQL for persistence.

## Quick Start

### Run with Docker Compose

1. Ensure Docker and Docker Compose are installed.
2. Build and run the application:
   ```sh
   make build-run
   ```
   This will start both the API and a PostgreSQL database.

> The app uses default environment variables. You can override them by creating a `.env` file if needed.

## Running Tests

Run all tests locally:

```sh
make run-tests
```

## Project Structure

This project follows the hexagonal architecture for clean separation of concerns:

```
.
├── cmd/            # Main application entry point
├── config/         # Configuration
├── internal/
│   ├── app/        # Core application logic
│   ├── httpjson/   # HTTP handler adapters
│   └── repository/ # Repository adapters
├── migrations/     # Database migration files and logic
```

## API Endpoints

| Method | Endpoint               | Request Body Example   | Response Example           | Description                                            |
| ------ | ---------------------- | ---------------------- | -------------------------- | ------------------------------------------------------ |
| GET    | /healthcheck           | -                      | 200 OK                     | Health check (returns 200 OK)                          |
| GET    | /packs/sizes           | -                      | [5, 10, 20]                | Returns a list of all available pack sizes             |
| POST   | /packs                 | {"size": 10}           | {"id": "uuid", "size": 10} | Create a new pack                                      |
| POST   | /packs/rebuild         | {"sizes": [5, 10, 20]} | [5, 10, 20]                | Rebuilds the pack list with the given sizes            |
| POST   | /packs/calculate-order | {"items": 23}          | {"10": 2, "3": 1}          | Calculates the minimum items/packs needed for an order |
