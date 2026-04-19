# Implementation Plan: FD Credit Score

## Goal

Build a working FD Credit Score prototype — a Go backend with scoring engine + minimal but visually impressive frontend — deployable as a single binary.

---

## Architecture Decision

**Go monolith** that serves both API and frontend:

```
┌──────────────────────────────────────────┐
│           Single Go Binary               │
│                                          │
│  ┌──────────────────────────────────┐    │
│  │  HTTP Server (net/http)          │    │
│  │                                  │    │
│  │  /api/score    POST → Scoring    │    │
│  │  /api/personas GET  → Mock Data  │    │
│  │  /api/insights POST → LLM       │    │
│  │  /*            GET  → Static     │    │
│  │                       Frontend   │    │
│  └──────────────────────────────────┘    │
│                                          │
│  ┌────────────┐  ┌───────────────────┐   │
│  │  Scoring   │  │  LLM Adapter      │   │
│  │  Engine    │  │  (OpenAI/Gemini)  │   │
│  └────────────┘  └───────────────────┘   │
│                                          │
│  ┌────────────────────────────────────┐  │
│  │  Mock Data (embedded JSON)         │  │
│  └────────────────────────────────────┘  │
└──────────────────────────────────────────┘
```

**Why this approach**:
- Single `go build` → single binary → deploy anywhere
- No CORS headaches (frontend and API on same origin)
- Go embeds static files with `embed` package — no separate file server needed
- Deploy on Render/Railway free tier in 2 minutes
- Shows Go engineering chops to judges

---

## User Review Required

> [!IMPORTANT]
> **LLM API choice**: We need an API key for generating natural language insights (score explanations, improvement suggestions). Options:
> 1. **OpenAI (GPT-4o-mini)** — cheapest, fastest, most reliable
> 2. **Google Gemini** — free tier available
> 3. **No LLM** — hardcode insight templates (faster to build, zero API cost, still works well)
> 
> **Recommendation**: Option 3 (hardcoded templates) for the qualifying round. The scoring engine IS the innovation — LLM-generated text is nice-to-have but not the core value. We can add LLM in the build phase if shortlisted.

> [!IMPORTANT]  
> **Deployment platform**: Render free tier (easiest for Go). Need a Render account. Alternative: Railway, Fly.io. Do you have a preference or existing account?

---

## Project Structure

```
/home/ayanami/meridian/sysc/
├── main.go                     # Entry point, HTTP server setup
├── go.mod
├── go.sum
│
├── internal/
│   ├── scoring/
│   │   ├── engine.go           # Core scoring algorithm (5 components)
│   │   ├── models.go           # Data types: DepositHistory, Score, etc.
│   │   ├── consistency.go      # Component 1: Consistency scorer
│   │   ├── discipline.go       # Component 2: Maturity discipline scorer
│   │   ├── growth.go           # Component 3: Growth trajectory scorer
│   │   ├── diversification.go  # Component 4: Diversification scorer
│   │   ├── intelligence.go     # Component 5: Tenure intelligence scorer
│   │   └── engine_test.go      # Unit tests for scoring
│   │
│   ├── insights/
│   │   └── generator.go        # Template-based insight generation
│   │
│   ├── personas/
│   │   └── personas.go         # 5 pre-built personas with embed data
│   │
│   └── handlers/
│       └── api.go              # HTTP handlers for all endpoints
│
├── web/                        # Frontend (embedded into Go binary)
│   ├── index.html              # Single page application
│   ├── css/
│   │   └── style.css           # All styling (dark theme, animations)
│   └── js/
│       ├── app.js              # Main app logic, routing
│       ├── scoring.js          # API calls, score display logic
│       └── charts.js           # Chart.js radar chart, score dial
│
├── data/
│   └── personas.json           # 5 persona deposit histories
│
└── README.md                   # Project documentation for GitHub
```

---

## Proposed Changes

### Backend — Core Scoring Engine

#### [NEW] [go.mod](file:///home/ayanami/meridian/sysc/go.mod)
Go module initialization. No external dependencies beyond stdlib required for core. Optional: `github.com/go-chi/chi` for cleaner routing (but `net/http` is fine for this scope).

#### [NEW] [models.go](file:///home/ayanami/meridian/sysc/internal/scoring/models.go)
Core data types:

```go
// Input: A user's deposit history
type DepositHistory struct {
    UserID    string       `json:"user_id"`
    Name      string       `json:"name"`
    Age       int          `json:"age"`
    City      string       `json:"city"`
    Deposits  []Deposit    `json:"deposits"`
}

type Deposit struct {
    ID              string    `json:"id"`
    Type            string    `json:"type"`              // "FD" or "RD"
    Bank            string    `json:"bank"`
    Amount          float64   `json:"amount"`            // Principal
    TenureMonths    int       `json:"tenure_months"`
    InterestRate    float64   `json:"interest_rate"`
    StartDate       string    `json:"start_date"`        // YYYY-MM-DD
    MaturityDate    string    `json:"maturity_date"`
    WithdrawnDate   *string   `json:"withdrawn_date"`    // null if held to maturity
    Status          string    `json:"status"`            // "active", "matured", "withdrawn"
    // RD-specific
    RDInstallment   float64   `json:"rd_installment"`    // monthly amount
    RDPaidMonths    int       `json:"rd_paid_months"`    // installments completed
    RDMissedMonths  int       `json:"rd_missed_months"`  // installments missed
}

// Output: The generated score
type ScoreResult struct {
    TotalScore       int                `json:"total_score"`       // 0-900
    ScoreBand        string             `json:"score_band"`        // "Excellent", "Good", etc.
    Components       []ComponentScore   `json:"components"`
    CreditProducts   []CreditProduct    `json:"credit_products"`
    Improvements     []Improvement      `json:"improvements"`
    Insights         []string           `json:"insights"`
    PatternDetected  string             `json:"pattern_detected"`  // e.g., "LIQUIDITY_GAP_SAVER"
    PeerPercentile   int                `json:"peer_percentile"`
}

type ComponentScore struct {
    Name        string  `json:"name"`
    Score       int     `json:"score"`       // 0-100
    MaxScore    int     `json:"max_score"`   // always 100
    Weight      float64 `json:"weight"`      // 0.30, 0.25, etc.
    Weighted    float64 `json:"weighted"`    // score * weight
    SubMetrics  []SubMetric `json:"sub_metrics"`
}
```

#### [NEW] [engine.go](file:///home/ayanami/meridian/sysc/internal/scoring/engine.go)
Orchestrator that calls all 5 component scorers, applies weights, detects patterns, and maps to 900 scale:

```go
func CalculateScore(history DepositHistory) ScoreResult {
    components := []ComponentScore{
        calculateConsistency(history),      // 30%
        calculateDiscipline(history),       // 25%
        calculateGrowth(history),           // 20%
        calculateDiversification(history),  // 15%
        calculateIntelligence(history),     // 10%
    }
    
    rawScore := sumWeighted(components)
    mappedScore := mapTo900(rawScore)
    pattern := detectPattern(components, history)
    products := recommendProducts(mappedScore, history, pattern)
    improvements := suggestImprovements(components, history)
    
    return ScoreResult{...}
}
```

#### [NEW] [consistency.go](file:///home/ayanami/meridian/sysc/internal/scoring/consistency.go)
- Deposit frequency: count FDs per year, normalize to 0-40
- Gap analysis: average months between consecutive FD start dates, map to 0-30
- Streak tracking: longest unbroken annual streak, map to 0-30

#### [NEW] [discipline.go](file:///home/ayanami/meridian/sysc/internal/scoring/discipline.go)
- FD completion rate: (matured / total closed) × 40
- Premature withdrawal penalty: deduct based on count and severity
- Average hold ratio: (actual hold / intended tenure) average, map to 0-30
- RD completion boost: completed RDs add bonus points (EMI readiness signal)

#### [NEW] [growth.go](file:///home/ayanami/meridian/sysc/internal/scoring/growth.go)
- YoY deposit growth: compare annual total deposits
- Recovery pattern detection: if growth dipped then recovered, give partial credit
- Corpus size: total active deposits, map to percentile bracket

#### [NEW] [diversification.go](file:///home/ayanami/meridian/sysc/internal/scoring/diversification.go)
- Bank count: unique banks used, map to 0-35
- Tenure spread: coefficient of variation of chosen tenures, map to 0-35
- Product mix: FD + RD + tax-saver, map to 0-30

#### [NEW] [intelligence.go](file:///home/ayanami/meridian/sysc/internal/scoring/intelligence.go)
- Rate optimization: did user choose higher-rate banks when available (compare against rate table)
- Ladder detection: are FD maturity dates staggered?
- Tax-saver awareness: has user created 80C tax-saver FDs?

#### [NEW] [engine_test.go](file:///home/ayanami/meridian/sysc/internal/scoring/engine_test.go)
Unit tests using the 5 personas — assert expected score ranges.

---

### Backend — API Layer

#### [NEW] [api.go](file:///home/ayanami/meridian/sysc/internal/handlers/api.go)

Three endpoints:

```
POST /api/score
  Body: DepositHistory JSON
  Returns: ScoreResult JSON

GET /api/personas
  Returns: List of 5 pre-built personas with their deposit histories

GET /api/personas/{id}/score
  Returns: Pre-calculated score for a specific persona
```

#### [NEW] [main.go](file:///home/ayanami/meridian/sysc/main.go)

```go
//go:embed web/*
var webFS embed.FS

func main() {
    mux := http.NewServeMux()
    
    // API routes
    mux.HandleFunc("POST /api/score", handlers.CalculateScore)
    mux.HandleFunc("GET /api/personas", handlers.ListPersonas)
    mux.HandleFunc("GET /api/personas/{id}/score", handlers.GetPersonaScore)
    
    // Frontend (embedded static files)
    webRoot, _ := fs.Sub(webFS, "web")
    mux.Handle("/", http.FileServer(http.FS(webRoot)))
    
    port := os.Getenv("PORT")
    if port == "" { port = "8080" }
    
    log.Printf("FD Credit Score running on :%s", port)
    http.ListenAndServe(":"+port, mux)
}
```

---

### Backend — Insight Generator

#### [NEW] [generator.go](file:///home/ayanami/meridian/sysc/internal/insights/generator.go)

Template-based (no LLM needed):

```go
// Pattern detection rules
func DetectPattern(components []ComponentScore, history DepositHistory) string {
    consistency := findComponent("Consistency", components)
    discipline := findComponent("Maturity Discipline", components)
    
    // High consistency + Low discipline = breaking FDs out of necessity
    if consistency.Score > 65 && discipline.Score < 55 {
        return "LIQUIDITY_GAP_SAVER"
    }
    // High everything except diversification
    if avgScore(components) > 70 && findComponent("Diversification", components).Score < 35 {
        return "LOYAL_SINGLE_BANK"
    }
    // Low consistency but recovery pattern in growth
    if consistency.Score < 40 && findComponent("Growth", components).Score > 50 {
        return "RECOVERING_SAVER"
    }
    // High across the board
    if minScore(components) > 70 {
        return "DISCIPLINED_OPTIMIZER"
    }
    return "STANDARD"
}

// Each pattern has a narrative template
var patternNarratives = map[string]string{
    "LIQUIDITY_GAP_SAVER": "You're a consistent saver who occasionally needs to break FDs for cash flow. A credit line could eliminate this — you'd stop losing FD interest and build credit simultaneously.",
    "LOYAL_SINGLE_BANK": "Your savings discipline is excellent, but keeping everything at one bank means you might be missing better rates elsewhere and exceeding DICGC insurance limits.",
    // ...
}
```

---

### Backend — Mock Data

#### [NEW] [personas.json](file:///home/ayanami/meridian/sysc/data/personas.json)

5 complete personas with realistic deposit histories (15-25 deposit records each). Each persona maps to the scenarios from the deep-dive:

1. **Priya** — Teacher, 9 years of annual FDs, single bank, perfect discipline
2. **Vikram** — Engineer, multi-bank optimizer, aggressive saver
3. **Ramesh** — Shopkeeper, COVID-disrupted, recovering
4. **Anita** — Freelancer, consistent creator but frequent breaker
5. **Sunita** — Retiree, 15-year single-bank ultra-conservative

---

### Frontend — Single Page Application

The frontend is minimal in code but visually impressive. Single HTML file with supporting CSS/JS.

#### [NEW] [index.html](file:///home/ayanami/meridian/sysc/web/index.html)

Single page with 4 views (no routing framework, just show/hide sections):

1. **Landing** → Hero with tagline + persona selector cards
2. **Score Reveal** → Animated score dial + radar chart + score band
3. **Breakdown** → Component-by-component detail with sub-metrics
4. **Credit Products** → What's unlocked + improvement roadmap

#### [NEW] [style.css](file:///home/ayanami/meridian/sysc/web/css/style.css)

Design system:
- **Dark theme** with deep navy background (`#0a0e1a`) 
- **Accent gradient**: teal-to-blue (`#00d4aa → #0066ff`) for score elements
- **Score dial**: Animated SVG circle that fills based on score (like a CIBIL reveal)
- **Cards**: Glassmorphism with subtle backdrop-blur
- **Typography**: Inter (Google Fonts) for clean, modern feel
- **Micro-animations**: Score counter animating from 0 → final score, radar chart drawing on
- **Responsive**: Works on laptop screen (demo mode) and mobile

Key visual elements:
```
┌─────────────────────────────────────────────────┐
│  FD CREDIT SCORE                                │
│                                                 │
│  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐      │
│  │Priya│ │Vikrm│ │Ramsh│ │Anita│ │Sunita│      │
│  │ 👩‍🏫 │ │ 👨‍💻 │ │ 🧑‍🔧 │ │ 👩‍🎨 │ │ 👵  │      │
│  └─────┘ └─────┘ └─────┘ └─────┘ └─────┘      │
│           [ + Custom Input ]                    │
│                                                 │
│      ┌─────────────────┐   ┌──────────────┐    │
│      │    ╭───────╮    │   │  Radar Chart │    │
│      │   ╱  689   ╲   │   │    ╱╲         │    │
│      │  │  / 900   │   │   │   /  \        │    │
│      │   ╲  GOOD  ╱   │   │  /____\       │    │
│      │    ╰───────╯    │   │              │    │
│      │   Score Dial    │   │              │    │
│      └─────────────────┘   └──────────────┘    │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │ ✅ FD-backed Credit Card: ₹6.8L limit   │   │
│  │ ✅ UPI Credit: ₹50K limit               │   │
│  │ ⚠️ Unsecured Card: ₹30K (improve to    │   │
│  │    unlock higher)                        │   │
│  └─────────────────────────────────────────┘   │
│                                                 │
│  Pattern: "LOYAL SINGLE BANK SAVER"            │
│  "Your discipline is exceptional, but..."       │
│                                                 │
│  📈 Improvement Path:                           │
│  ┌─ Open FD at second bank ──── +35 pts ──┐   │
│  └─ Try a 6-month tenure ───── +15 pts ──┘   │
└─────────────────────────────────────────────────┘
```

#### [NEW] [app.js](file:///home/ayanami/meridian/sysc/web/js/app.js)
- Fetch personas from `/api/personas`
- On persona click → POST to `/api/score` → animate score reveal
- Custom input form: add deposits manually → submit for scoring
- View transitions with CSS animations

#### [NEW] [charts.js](file:///home/ayanami/meridian/sysc/web/js/charts.js)
- **Score dial**: SVG `<circle>` with `stroke-dasharray` animation (CSS-driven, no library needed)
- **Radar chart**: Chart.js (loaded via CDN) — 5-axis radar for component scores
- **Score counter**: JS counter animation from 0 → final score over 1.5 seconds

---

## Execution Timeline

### Phase 1: Backend Core (Hours 1-6)

| Hour | Task | Output |
|------|------|--------|
| 1 | Go project setup, `go mod init`, directory structure, models.go | Compilable project |
| 2 | Mock persona data — write `personas.json` with all 5 personas, complete deposit histories | 5 realistic personas |
| 3 | `consistency.go` + `discipline.go` — first 2 scoring components | 55% of scoring logic |
| 4 | `growth.go` + `diversification.go` + `intelligence.go` — remaining 3 components | 100% of scoring logic |
| 5 | `engine.go` — orchestrator, pattern detection, credit product recommendations | Complete scoring engine |
| 6 | `engine_test.go` — test all 5 personas, verify scores match expected ranges | Tested backend |

### Phase 2: API + Server (Hours 7-8)

| Hour | Task | Output |
|------|------|--------|
| 7 | `api.go` — HTTP handlers, JSON serialization | Working API |
| 8 | `main.go` — server setup, embed frontend, test with curl | API responding to requests |

### Phase 3: Frontend (Hours 9-14)

| Hour | Task | Output |
|------|------|--------|
| 9 | `index.html` — page structure, all sections | HTML skeleton |
| 10 | `style.css` — dark theme, glassmorphism cards, typography, layout | Styled but static |
| 11 | `style.css` continued — score dial SVG, animations, responsive | Visually complete |
| 12 | `app.js` — persona loading, API calls, view transitions | Interactive |
| 13 | `charts.js` — radar chart (Chart.js), score counter animation | Visualizations working |
| 14 | Custom input form — let user add deposits manually | Full feature set |

### Phase 4: Polish + Deploy (Hours 15-18)

| Hour | Task | Output |
|------|------|--------|
| 15 | Integration testing — full flow from persona click to score display | Working app |
| 16 | Visual polish — animations, transitions, mobile responsive | Polished app |
| 17 | Deploy to Render/Railway, verify live URL works | Live demo URL |
| 18 | Write README.md, record 2-min demo video, submit | **SUBMITTED** |

---

## Open Questions

> [!IMPORTANT]
> 1. **Do you have an OpenAI/Gemini API key**, or should we go fully template-based for insights? (Recommendation: templates first, LLM if time permits)
> 2. **Render.com account** — do you have one, or should we plan for Railway/Fly.io?
> 3. **The `sysc` directory** — is this where you want the project, or should we create a new subdirectory like `sysc/fd-credit-score/`?

---

## Verification Plan

### Automated Tests
```bash
# Run scoring engine unit tests
go test ./internal/scoring/ -v

# Test API endpoints
curl -X POST http://localhost:8080/api/score -d @data/personas.json
curl http://localhost:8080/api/personas

# Build and verify binary runs
go build -o fd-credit-score .
./fd-credit-score
```

### Manual Verification
- Load each persona in browser → verify score matches expected range
- Test custom input flow → verify score calculation is reasonable
- Test on mobile viewport → verify responsive layout
- Test deployed URL → verify it works end-to-end

### Browser Testing
- Use browser tool to navigate the live demo
- Click through all 5 personas
- Verify animations, charts, and transitions work smoothly
- Record a demo video for submission
