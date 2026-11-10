# Konkur Pulse 

Konkur Pulse is a high-performance, resilient backend API built in Go. It acts as a study allocation engine and session management system for students, integrating with external exam simulators (like Maze) to dynamically generate daily study plans based on real-time academic performance.

This project demonstrates advanced system design principles, including high-concurrency handling, upstream resilience, idempotency, and observability.

##  Key Features & Architecture

* **Resilient External Communication:** Implements a custom HTTP client with an Exponential Backoff & Retry mechanism to gracefully handle rate limits (429) and timeouts from unstable upstream APIs.

* **High-Performance Caching:** Utilizes a thread-safe `sync.Map` in-memory cache to bypass upstream network bottlenecks, preventing OS socket exhaustion (TCP port limits) and ensuring sub-millisecond response times under heavy load.

* **Idempotent Operations:** Prevents duplicate study session entries caused by mobile client retries or poor network conditions by enforcing unique `client_session_id` constraints at the PostgreSQL database level (`ON CONFLICT DO NOTHING`).

* **Observability & Monitoring:** Fully instrumented with structured JSON logging (`slog`), unique OpenTelemetry-style Request IDs for distributed tracing, and Prometheus metrics for real-time monitoring of RPS, P95/P99 latencies, and upstream error rates.

* **Standard Library Routing:** Leverages the native Go 1.22+ `net/http` ServeMux for clean, dependency-free HTTP routing.

## 📂 Project Structure

```
konkur-pulse/
├── cmd/
│   ├── api/          # Main application entrypoint
│   └── mock/         # Mock Maze simulator (for testing upstream instability)
├── internal/
│   ├── handler/      # HTTP handlers (plan generation, session saving)
│   ├── metrics/      # Prometheus metrics definitions
│   ├── middleware/   # Observability and Request ID injection
│   ├── repository/   # PostgreSQL DB operations and external API clients
│   └── service/      # Core business logic (study time allocation engine)
├── docs/adr/         # Architecture Decision Records (ADRs 001 to 005)
├── docker-compose.yml# PostgreSQL database container configuration
└── go.mod            # Go module dependencies

```

## 🛠️ Tech Stack

* **Language:** Go (Golang)

* **Database:** PostgreSQL (with `lib/pq` driver)

* **Monitoring:** Prometheus (`client_golang`)

* **Logging:** standard `log/slog`

* **Infrastructure:** Docker (for DB), WSL/Linux (for optimal high-concurrency deployment)

##  Getting Started

### 1. Prerequisites

* Go 1.22 or higher

* Docker and Docker Compose

* Linux/WSL (Recommended for high-concurrency load testing)

### 2. Setup Database

Start the PostgreSQL container:

```
docker-compose up -d

```

### 3. Run the Servers

Run the Mock Simulator in the background (runs on port 8081):

```
go run ./cmd/mock/main.go &

```

Run the Main API Server (runs on port 8080):

```
go run ./cmd/api/main.go

```

## 📡 API Endpoints

### 1. Get Study Plan

Fetches exam results from the upstream simulator and allocates study minutes per subject.

* **GET** `/students/{id}/plan?date=YYYY-MM-DD`

* **Response (200 OK):**

```
{
  "plan": {
    "ریاضی": 75,
    "شیمی": 120,
    "فیزیک": 96
  },
  "student_id": "123"
}

```

### 2. Save Study Session

Idempotent endpoint to record a completed study session.

* **POST** `/students/{id}/sessions`

* **Payload:**

```
{
  "date": "2026-10-30",
  "subject": "فیزیک",
  "duration_minutes": 90,
  "client_session_id": "unique-uuid-from-mobile"
}

```

### 3. Metrics & Health

* **GET** `/health` (API health check)

* **GET** `/metrics` (Prometheus metrics exposition)

##  Load Testing & Performance

The system was aggressively load-tested using [Vegeta](https://github.com/tsenart/vegeta). Thanks to the in-memory cache and connection multiplexing, the Go server successfully handled **10,000 Requests Per Second (RPS)** without dropping a single connection when deployed in a Linux environment (bypassing Windows TCP limits).

To reproduce the benchmark in WSL/Linux:

```
ulimit -n 200000
echo "GET http://127.0.0.1:8080/students/123/plan?date=2026-10-30" | vegeta attack -duration=20s -rate=10000 -connections=500 -keepalive=true | vegeta report

```

## 👤 Author

Developed by **Matin Hajidavalou**