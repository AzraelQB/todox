# Todox API

## Build

To build the api application, you need to have Golang installed on your machine. Clone the repository and run the following commands:

```bash
go build -o app
```

# Test
To run tests, use the following command:

```bash
go test ./...
```

# Run in Docker
Make sure you have Docker installed on your machine. Create a Docker image and run the todox api using the provided Dockerfile and Docker Compose file:

```bash
docker-compose up --build
```

This command will build the Docker image and start the services defined in the docker-compose.yml file. Todox api will be accessible at http://localhost:8080.

# API Documentation

Generate swagger docs
```bash
go install github.com/swaggo/swag/cmd/swag@latest

swag init
```
Swagger documentation is available at http://localhost:8080/swagger/index.html. You can use this documentation to explore and test your API endpoints.
