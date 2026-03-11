# Zguba.gov

Lost and found item registration system for Polish public administration offices. Provides a multi-step wizard for municipal clerks to register found items, a public search interface, and a standards-compliant API for integration with [dane.gov.pl](https://dane.gov.pl).

Built in 24 hours during the [HackNation](https://hacknation.pl/#about) hackathon.

## Architecture

The application is a Go binary with embedded static assets, backed by a MySQL database. It serves both the HTML interface (rendered server-side with HTMX) and the REST/OData API.

```
    ┌───────────────────────────────────────────────────┐
    │                  Zguba.gov                        │
    │                                                   │
    │   ┌─────────────┐    ┌──────────────────────┐     │
    │   │ Gin Router  │───>│ Page Handlers        │     │
    │   │             │    │ (HTMX partials)      │     │
    │   │ /           │    └──────────┬───────────┘     │
    │   │ /steps/*    │               │                 │
    │   │ /api/*      │    ┌──────────▼───────────┐     │
    │   │ /odata/*    │───>│ API / OData Handlers │     │
    │   │ /metadata   │    └──────────┬───────────┘     │
    │   │ /health     │               │                 │
    │   └─────────────┘    ┌──────────▼───────────┐     │
    │                      │ Repository Layer     │     │
    │   ┌─────────────┐    └──────────┬───────────┘     │
    │   │ Embedded    │               │                 │
    │   │ Templates   │    ┌──────────▼───────────┐     │
    │   │ & Static    │    │ MySQL Database       │     │
    │   │ Assets      │    │ (zguba_gov)          │     │
    │   └─────────────┘    └──────────────────────┘     │
    │                                                   │
    └───────────────────────────────────────────────────┘
```

The request flow for the wizard UI:

1. **Browser loads the page** — Gin renders the full layout with the first wizard step.
2. **User fills in a step** — HTMX sends a POST to `/steps/next` with all form data.
3. **Server validates and responds** — returns the next step partial (or an error modal). Hidden form fields preserve state across steps.
4. **Final submission** — the completed form is POSTed to `/api/found-items`. On success, a confirmation modal is shown and the items list refreshes.

## Key features

- Multi-step wizard for registering found items with server-side validation
- Territorial unit autocomplete covering all Polish voivodeships, counties, and municipalities (2809 units in MySQL)
- Items listing with filtering by category, municipality, status, and free-text search
- JSON and CSV export of found items
- RESTful API with full CRUD operations
- OData-compatible endpoint with `$filter`, `$orderby`, `$top`, `$skip`, `$count`
- DCAT-AP metadata endpoint for dane.gov.pl catalog integration
- Responsive UI following GOV.PL design guidelines
- Health check endpoint for container orchestration

## Getting started

### Prerequisites

- Go 1.24+ — [download](https://go.dev/dl/)
- MySQL 8+ — [download](https://dev.mysql.com/downloads/)
- Docker (recommended)

### Run with Docker

```bash
docker compose up --build
```

The application starts on [http://localhost](http://localhost). MySQL is initialized automatically with the schema and territorial units data from `init.sql`.

### Run locally

1. Start a MySQL instance and create the database:

```bash
mysql -u root -p < init.sql
```

2. Run the application:

```bash
DATABASE_URL="root:password@tcp(localhost:3306)/zguba_gov?parseTime=true&charset=utf8mb4" go run ./cmd/server
```

The application starts on [http://localhost:80](http://localhost:80).

### Build

```bash
go build -o zguba-gov ./cmd/server
```

## Configuration

### Optional environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `80` | HTTP server port |
| `DATABASE_URL` | `root:password@tcp(localhost:3306)/zguba_gov?parseTime=true&charset=utf8mb4` | MySQL DSN |
| `CORS_ORIGINS` | `http://localhost:4200,http://localhost:3000` | Comma-separated list of allowed CORS origins |

## API endpoints

### Found items

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/found-items` | List items (query: `skip`, `limit`, `category`, `municipality`, `status`, `search`) |
| `POST` | `/api/found-items` | Create item |
| `GET` | `/api/found-items/:id` | Get item by ID |
| `PUT` | `/api/found-items/:id` | Update item |
| `DELETE` | `/api/found-items/:id` | Delete item |
| `GET` | `/api/found-items/categories/list` | List available categories |
| `GET` | `/api/stats` | Get statistics |

### OData

| Method | Path | Description |
|---|---|---|
| `GET` | `/odata/FoundItems` | Query items (`$filter`, `$orderby`, `$top`, `$skip`, `$count`) |
| `GET` | `/odata/$metadata` | EDMX metadata document |

### Metadata (DCAT-AP)

| Method | Path | Description |
|---|---|---|
| `GET` | `/metadata` | Dataset catalog |
| `GET` | `/metadata/distribution/:id` | Distribution info |

### Health

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Health check |

## Troubleshooting

| Problem | Possible cause | Resolution |
|---|---|---|
| Autocomplete shows no results | Query too short | Type at least 2 characters to trigger autocomplete |
| Connection refused to MySQL | MySQL not running or wrong DSN | Check that MySQL is running and `DATABASE_URL` is correct |
| Port already in use | Another process on port 80 | Set `PORT` environment variable to a different port |
| Static assets not loading | Modified embedded files without rebuilding | Run `go build` again — assets are embedded at compile time |
| `init.sql` not loaded | MySQL volume already initialized | Remove the `mysql-data` volume (`docker compose down -v`) and restart |

## Security considerations

- The application does not implement authentication — it is designed for internal use within municipal offices
- Input is validated server-side on every wizard step before database insertion
- SQL queries use parameterized statements to prevent SQL injection
- CORS origins are configurable and restricted by default
- The Docker image uses a minimal Alpine base

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.
