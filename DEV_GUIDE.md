# Developer Guide

## Project Structure

```
├── main.go                          # Entry point, HTTP server setup
├── go.mod                           # Go module definition
│
├── internal/
│   ├── scoring/
│   │   ├── models.go                # DepositHistory, ScoreResult, ComponentScore types
│   │   ├── engine.go                # CalculateScore() - orchestrator
│   │   ├── consistency.go           # calculateConsistency() - 30% weight
│   │   ├── discipline.go            # calculateDiscipline() - 25% weight
│   │   ├── growth.go                # calculateGrowth() - 20% weight
│   │   ├── diversification.go      # calculateDiversification() - 15% weight
│   │   └── intelligence.go         # calculateIntelligence() - 10% weight
│   │
│   ├── handlers/
│   │   └── api.go                   # HTTP handlers for /api/* endpoints
│   │
│   ├── personas/
│   │   ├── personas.go              # GetAll(), GetByID() - loads embedded JSON
│   │   └── data/personas.json       # 5 pre-built persona profiles
│   │
│   └── insights/
│       └── generator.go             # Template-based pattern narratives
│
└── web/                             # Static frontend files
    ├── index.html                   # SPA with 3 views (landing, custom, score)
    ├── css/style.css                # Dark theme, glassmorphism, animations
    └── js/
        ├── app.js                   # API calls, view transitions, form handling
        └── charts.js                # Chart.js radar chart setup
```

## Key Functions

### Scoring Engine (internal/scoring/)

**engine.go:CalculateScore(history DepositHistory) ScoreResult**
- Main orchestrator
- Calls all 5 component calculators
- Maps raw score to 900 scale
- Detects behavioral pattern
- Recommends credit products
- Suggests improvements
- Generates insights

**Component scorers** (each returns ComponentScore):
- `calculateConsistency()` - FD frequency, gap analysis, streak
- `calculateDiscipline()` - Completion rate, early withdrawal penalty, hold ratio
- `calculateGrowth()` - YoY growth, recovery patterns, corpus size
- `calculateDiversification()` - Bank count, tenure spread, product mix
- `calculateIntelligence()` - Rate optimization, ladder detection, tax-saver awareness

### API Handlers (internal/handlers/)

- `CalculateScore(w, r)` - POST /api/score
- `ListPersonas(w, r)` - GET /api/personas
- `GetPersonaScore(w, r)` - GET /api/personas/{id}/score
- `GetHealth(w, r)` - GET /api/health

### Persona Loader (internal/personas/)

- Uses Go's `//go:embed` to embed personas.json into binary
- `GetAll()` returns []Persona
- `GetByID(id)` returns *Persona

## Testing

### Build and Run
```bash
go build -o fd-credit-score .
./fd-credit-score
# Server on http://localhost:8080
```

### Test API Endpoints
```bash
# Health check
curl http://localhost:8080/api/health

# List personas
curl http://localhost:8080/api/personas

# Get specific persona
curl http://localhost:8080/api/personas/priya

# Calculate score for a persona
curl http://localhost:8080/api/personas/priya/score

# Custom score calculation
curl -X POST http://localhost:8080/api/score \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test",
    "name": "Test User",
    "age": 30,
    "city": "Delhi",
    "deposits": [
      {
        "type": "FD",
        "bank": "SBI",
        "amount": 100000,
        "tenure_months": 12,
        "interest_rate": 6.5,
        "start_date": "2024-01-01",
        "maturity_date": "2025-01-01",
        "status": "active"
      }
    ]
  }'
```

### Test Frontend
```bash
# Open in browser
http://localhost:8080

# Should see:
# - 5 persona cards on landing page
# - Click any to see score breakdown
# - "Custom Input" button for manual entry
```

### Run Unit Tests
```bash
go test ./internal/scoring/ -v
```

## Adding New Personas

Edit `internal/personas/data/personas.json`:
```json
{
  "id": "unique-id",
  "name": "Person Name",
  "age": 35,
  "occupation": "Job Title",
  "city": "City",
  "deposits": [
    {
      "type": "FD",  // or "RD"
      "bank": "Bank Name",
      "amount": 100000,
      "tenure_months": 12,
      "interest_rate": 6.5,
      "start_date": "2024-01-01",
      "maturity_date": "2025-01-01",
      "status": "active",  // "active", "matured", "withdrawn"
      "withdrawn_date": null  // set if withdrawn early
    }
  ]
}
```

## Modifying Scoring Logic

Each component is in its own file under `internal/scoring/`:
- Weights are defined in `engine.go` (currently: 0.30, 0.25, 0.20, 0.15, 0.10)
- Score ranges are 0-100 per component, mapped to 0-900 total
- Pattern detection logic in `engine.go:detectPattern()`
- Product recommendations in `engine.go:recommendProducts()`

## Common Issues

**Build fails**: Check `go.mod` and ensure all imports resolve with `go mod tidy`

**API returns 404**: Verify routing in `main.go` - pattern must match exactly

**Frontend not loading**: Check `main.go` serves from `./web` directory

**Score seems off**: Review persona deposit data and component logic in respective files