# Dig Your Movie

A simple DNS server that returns movie descriptions based on OMDB data.

## Project Structure

- `cmd/server`: The DNS server entry point.
- `cmd/client`: A CLI tool to query the DNS server.
- `internal/config`: Configuration management.
- `internal/dns`: DNS server implementation.
- `internal/omdb`: OMDB API client.

## Getting Started

### Prerequisites

- Go 1.20+
- OMDB API Key (default provided for dev)

### Running the Server

```bash
go run cmd/server/main.go
```

The server listens on port 8095 by default.

### Running the Client

```bash
go run cmd/client/main.go "The Matrix"
```

## Docker

To build and run with Docker:

```bash
docker build -t dig-your-movie .
docker run -p 8095:8095 dig-your-movie
```
