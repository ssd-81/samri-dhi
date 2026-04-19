# FD Credit Score

A Go-based credit scoring system that evaluates Fixed Deposit (FD) holders based on their deposit history. The scoring engine analyzes 5 key components to generate a credit score (0-900) and provides personalized insights, credit product recommendations, and improvement suggestions.

## Architecture

Single Go binary that serves both API and frontend.

## Scoring Components

| Component | Weight |
|-----------|--------|
| Consistency | 30% |
| Maturity Discipline | 25% |
| Growth Trajectory | 20% |
| Diversification | 15% |
| Tenure Intelligence | 10% |

## Run

```bash
go build -o fd-credit-score .
./fd-credit-score
# Server runs on http://localhost:8080
```

## API Endpoints

- `POST /api/score` - Calculate score
- `GET /api/personas` - List personas
- `GET /api/personas/{id}/score` - Get persona score