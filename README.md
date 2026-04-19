# Samridhi: The FD-Backed Credit Score 📈

> *Turning your savings discipline into your financial superpower.*

## The Problem: "Credit Invisible, but Cash Rich"
Millions of people consistently save money in Fixed Deposits (FDs) year after year. They demonstrate immense financial discipline, patience, and wealth-building capability. Yet, when they apply for a premium credit card or a low-interest loan, they are often rejected or given a low score by traditional credit bureaus (like CIBIL). 

Why? Because traditional scoring models only measure how well you *borrow* money, not how well you *save* it.

## The Solution: A New Dimension of Creditworthiness
Samridhi is a novel credit scoring engine that evaluates fixed deposit history to generate a comprehensive credit score (300-900). By analyzing saving habits, we can prove to lenders that a user is a low-risk, high-discipline individual deserving of premium credit products.

We look beyond generic credit history and dive deep into actual savings behavior:
- **Consistency (30%):** Do you save regularly every year? Have you built a solid streak?
- **Maturity Discipline (25%):** Do you let your FDs mature, or do you frequently break them early for liquidity?
- **Growth Trajectory (20%):** Is your overall corpus growing? Do you recover well after a financial dip?
- **Diversification (15%):** Do you spread your risk across different banks and deposit types (FDs, RDs)?
- **Tenure Intelligence (10%):** Do you optimize for better interest rates and utilize tax-saver FDs?

## The Output
Based on this analysis, the engine provides:
1. **A Custom Score (300-900):** Easily understood by both users and financial institutions.
2. **Behavioral Personas:** Identifies user saving styles (e.g., *Disciplined Optimizer*, *Liquidity Gap Saver*).
3. **Actionable Insights:** Suggests precisely what needs to be done to improve the score.
4. **Credit Recommendations:** Unlocks specific credit products based on the user's FD profile (like FD-backed cards or low-rate personal loans).

---

## 🛠️ Technical Overview

Samridhi is built as a lightning-fast, single-binary Go application that serves both a robust REST API and a beautiful, interactive frontend. 

### Architecture
- **Backend:** Go (1.25+) with zero external dependencies for the core scoring engine.
- **Frontend:** Vanilla HTML/CSS/JS (SPA) with Chart.js for data visualization. Modern, responsive Glassmorphism design system.
- **Data Model:** Embedded JSON personas for instant demonstration and testing.

### Running Locally

```bash
# Build the binary
go build -o fd-credit-score .

# Run the server
./fd-credit-score

# The application and API will be running on http://localhost:8080
```

### Core API Endpoints
- `POST /api/score` - Calculate a custom score based on provided deposit history
- `GET /api/personas` - List all predefined user personas
- `GET /api/personas/{id}/score` - Run the scoring engine on a specific embedded persona

*For detailed information on the scoring logic, weights, and component architecture, please refer to the [Developer Guide](DEV_GUIDE.md).*