# Go URL Shortener with Rate Limiting, Expiry & Analytics

A production-ready URL shortener written in Go. Features include:

- URL shortening with TTL (expiration support)
- Redirect endpoint with expiry validation
- Analytics tracking (hit count, timestamps)
- Per-IP rate limiting (10 req/min default)
- In-memory storage (PostgreSQL/Redis-ready)
- Docker-ready for deployment
- Modular folder structure

---

## Getting Started

### Prerequisites
- Go 1.20+
- (Optional) Docker

### Install Dependencies
```bash
go mod tidy
```

### Run the Server
```bash
go run ./cmd/server
```

Server starts on `http://localhost:8000`.

---

## API Endpoints

### POST `/shorten`

Shorten a long URL.

**Request:**
```json
{
  "url": "https://google.com",
  "ttl_seconds": 3600
}
```

**Response:**
```json
{
  "short_url": "http://localhost:8000/abc123"
}
```

---

### GET `/{code}`

Redirects to the original URL if not expired.

---

### GET `/analytics/{code}`

Returns analytics for a short URL.

**Response:**
```json
{
  "original_url": "https://google.com",
  "created_at": "2025-05-04T18:32:16Z",
  "expires_at": "2025-05-04T19:32:16Z",
  "access_count": 3,
  "last_access": "2025-05-04T18:40:12Z"
}
```

---

## Rate Limiting

The app enforces per-IP rate limiting:
- Max 10 requests per minute
- Excess requests receive HTTP `429 Too Many Requests`

---

## Project Structure

```
go-url-shortener/
├── cmd/server          # Entry point
├── internal/
│   ├── api/            # Handlers
│   ├── shortener/      # Core URL logic
│   ├── analytics/      # (optional analytics pkg)
│   ├── rateLimiter/    # Middleware
│   └── storage/        # DB abstraction (if needed)
├── pkg/                # Utils (logger, config, etc.)
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

---

## Docker Support (Optional)
```bash
docker build -t go-url-shortener .
docker run -p 8000:8000 go-url-shortener
```

---

## Roadmap

- [x] In-memory storage
- [x] Rate limiting
- [x] Analytics
- [ ] Background expiry cleanup
- [ ] Redis/PostgreSQL integration
- [ ] OpenTelemetry support

---

## Author

Built by [Aditya Padgal](https://github.com/adityapadgal)
