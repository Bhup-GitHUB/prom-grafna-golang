# Prometheus and Grafana with Golang

A simple Go application demonstrating Prometheus metrics collection and Grafana visualization.

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose (for Prometheus and Grafana)

## Installation

1. Install dependencies:
```bash
go mod download
```

2. Run the application:
```bash
go run main.go
```

The server will start on port 3000 (or the port specified in the PORT environment variable).

## Endpoints

- `GET /health` - Health check endpoint
- `GET /user` - User endpoint
- `GET /math` - Math endpoint (simulates heavy computation)
- `GET /metrics` - Prometheus metrics endpoint

## Prometheus Setup

1. Start Prometheus and Grafana using Docker Compose:
```bash
docker-compose up -d
```

2. Prometheus will be available at: http://localhost:9090
3. Grafana will be available at: http://localhost:3001 (default credentials: admin/admin)

## Configure Prometheus

Prometheus is pre-configured to scrape metrics from the Go application at `http://localhost:3000/metrics`.

## Configure Grafana

1. Login to Grafana at http://localhost:3001
2. Add Prometheus as a data source:
   - URL: http://prometheus:9090
   - Access: Server (default)
3. Import dashboards or create your own queries using the following metrics:
   - `http_requests_total` - Total HTTP requests
   - `http_active_requests` - Active HTTP requests
   - `http_request_duration_seconds` - Request duration histogram

## Metrics

The application exposes three types of Prometheus metrics:

1. **Counter**: `http_requests_total` - Tracks total number of HTTP requests
2. **Gauge**: `http_active_requests` - Tracks current number of active requests
3. **Histogram**: `http_request_duration_seconds` - Tracks request duration distribution

## Testing

Generate some traffic to see metrics:
```bash
curl http://localhost:3000/health
curl http://localhost:3000/user
curl http://localhost:3000/math
```

View metrics:
```bash
curl http://localhost:3000/metrics
```

