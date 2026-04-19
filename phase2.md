# Implementation Plan: FD Credit Score → Full Experience Redesign

## Goal

Transform the FD Credit Score from a scoring calculator into a storytelling product that makes hackathon judges *feel* the problem and the solution. Backend stays unchanged except one small addition. All work is frontend-heavy.

---

## Architecture Overview

```
CURRENT FLOW:
  Landing (persona cards) ──click──→ Score Page (number + radar + lists)

NEW FLOW:
  Landing (rich persona cards with backstory)
    │
    ├── click persona ──→ Profile Page (story + context + "Analyze" CTA)
    │                          │
    │                          └──→ Score Reveal (ceremony + full breakdown)
    │                                    │
    │                                    ├── Component Breakdown (expandable)
    │                                    ├── Pattern Narrative (full story)
    │                                    ├── Credit Products (premium cards)
    │                                    ├── Improvement Path (with projections)
    │                                    └── Deposit Timeline (chart)
    │
    └── "Your Own Score" ──→ Form Wizard (step-by-step deposit entry)
                                  │
                                  └──→ Score Reveal (same as above)
```

---

## Backend Change: One New Field in API Response

#### [MODIFY] [api.go](file:///home/ayanami/workspace/github.com/ssd-81/samri-dhi/internal/handlers/api.go)

Add persona metadata to the score response so the frontend can show backstory without a separate API call.

When hitting `GET /api/personas/{id}/score`, the response should include:

```json
{
  "persona": {
    "name": "Priya Sharma",
    "age": 38,
    "occupation": "Teacher", 
    "city": "Jaipur",
    "deposit_count": 11,
    "total_corpus": 1440000,
    "years_active": 9,
    "active_deposits": 2
  },
  "total_score": 689,
  "score_band": "Good",
  "components": [...],
  "credit_products": [...],
  "improvements": [...],
  "insights": [...],
  "pattern_detected": "LOYAL_SINGLE_BANK",
  "peer_percentile": 75,
  "projected_score": 739,
  "cibil_equivalent": "700-749"
}
```

New fields to compute in `engine.go`:
- `projected_score`: Apply the top improvement suggestion and recalculate
- `cibil_equivalent`: Simple mapping from FD score bands to approximate CIBIL ranges
- Persona summary stats: computed from deposit history (count, total corpus, years active)

> [!NOTE]
> This is a small backend change — add a wrapper struct that includes persona info + score result. ~30 minutes of work.

---

## Screen 1: Landing Page (Redesigned)

### What Changes
The current landing has a hero title + plain persona cards. The new version makes each persona card tell a micro-story and adds a problem statement section.

### Layout

```
┌─────────────────────────────────────────────────────────────┐
│  ◈ Samridhi                                        [About]  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│           500 Million Indians Save Diligently.              │
│           Zero Get Credit For It.                           │
│                                                             │
│     Traditional credit scores only measure debt repayment.  │
│     Samridhi measures what matters: savings discipline.     │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  ₹66L Cr          400-500M           0              │   │
│  │  FD Holdings      Credit Invisible   Credit Score   │   │
│  │  in India         Indians            for savers     │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                             │
│     ─────── See how it works ───────                        │
│                                                             │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌────────────────┐ │
│  │ 👩‍🏫      │  │ 👨‍💻      │  │ 🧑‍🔧      │  │ + Calculate    │ │
│  │ Priya   │  │ Vikram  │  │ Ramesh  │  │   Your Own     │ │
│  │ Teacher │  │ Engineer│  │ Shop    │  │   Score        │ │
│  │ 38, Jai │  │ 32, Blr │  │ 45, Del │  │                │ │
│  │         │  │         │  │         │  │   Enter your   │ │
│  │ 9 yrs   │  │ 15 FDs  │  │ COVID   │  │   FD history   │ │
│  │ perfect │  │ 8 banks │  │ recover │  │   and discover  │ │
│  │ saving  │  │ optimzd │  │ story   │  │   your score   │ │
│  └─────────┘  └─────────┘  └─────────┘  └────────────────┘ │
│                                                             │
│  ┌─────────┐  ┌─────────┐                                   │
│  │ 👩‍🎨      │  │ 👵      │                                   │
│  │ Anita   │  │ Sunita  │                                   │
│  │ Freelncr│  │ Retiree │                                   │
│  │ 29, Mum │  │ 62, Chn │                                   │
│  │         │  │         │                                   │
│  │ Breaks  │  │ 15 year │                                   │
│  │ FDs for │  │ single  │                                   │
│  │ cashflw │  │ bank    │                                   │
│  └─────────┘  └─────────┘                                   │
│                                                             │
│     Built for the Blostem AI Builder Hackathon 2026         │
└─────────────────────────────────────────────────────────────┘
```

### Persona Card Design

Each card has:
- **Avatar**: Gradient circle with emoji (existing)
- **Name + Occupation**: Bold name, muted occupation (existing, enhanced)
- **Location badge**: "38, Jaipur" in a small pill
- **Tagline**: 2-line story hook unique to each persona (NEW — hardcoded in frontend)
- **Stat pill**: "11 FDs • 9 years" at the bottom (NEW — computed from persona data)
- **Hover effect**: Card lifts, border glows teal, subtle scale(1.02)

Persona taglines (hardcoded in JS):
```javascript
const personaTaglines = {
  priya:  "9 years of perfect saving. Zero credit history.",
  vikram: "8 banks, rate-optimized, aggressive saver.",
  ramesh: "COVID broke his streak. He's coming back stronger.",
  anita:  "She saves consistently... then breaks her FDs for cash flow.",
  sunita: "15 years at one bank. ₹13L in FDs. Still invisible to lenders."
};
```

### Stats Bar (NEW component)

Three number cards in a row with count-up animation on page load:
- ₹66 Lakh Crore — "FD Holdings in India"
- 400-500 Million — "Credit Invisible Indians"  
- 0 — "Credit Scores for Savers"

These animate in with a staggered fade-up (200ms delay between each).

### "Calculate Your Own Score" Card

This is the 6th card in the grid, but styled differently:
- Dashed border instead of solid
- Plus icon + "Calculate Your Own Score" as title
- Subtitle: "Enter your FD history and discover your score"
- Gradient text treatment on "Your Own Score"

---

## Screen 2: Profile Page (NEW — The Story Beat)

### Purpose
This is the "pause before the reveal" — it shows WHO this person is and WHY their score matters. This is the emotional beat that gives the score meaning.

### When It Appears
After clicking a persona card, BEFORE the score loads.

### Layout

```
┌─────────────────────────────────────────────────────────────┐
│  ◈ Samridhi                                    [← Back]     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│     ┌────────────────────────────────────────────────┐      │
│     │                                                │      │
│     │   👩‍🏫                                           │      │
│     │                                                │      │
│     │   Priya Sharma                                 │      │
│     │   Teacher • 38 • Jaipur                        │      │
│     │                                                │      │
│     └────────────────────────────────────────────────┘      │
│                                                             │
│     ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌─────────┐  │
│     │   11     │  │   ₹14.4L │  │   9      │  │   2     │  │
│     │ Deposits │  │  Corpus  │  │  Years   │  │ Active  │  │
│     └──────────┘  └──────────┘  └──────────┘  └─────────┘  │
│                                                             │
│     ┌────────────────────────────────────────────────┐      │
│     │  "Priya has created an FD every single year    │      │
│     │   since 2016. She held every one to maturity.  │      │
│     │   She's completed two Recurring Deposits       │      │
│     │   without missing an installment.              │      │
│     │                                                │      │
│     │   Her CIBIL score? Not Available.              │      │
│     │   She has never taken a loan or credit card.   │      │
│     │                                                │      │
│     │   Traditional credit scoring can't see her.    │      │
│     │   Samridhi can."                               │      │
│     └────────────────────────────────────────────────┘      │
│                                                             │
│            ┌─────────────────────────────┐                  │
│            │  ✨ Analyze Her FD History  │                  │
│            └─────────────────────────────┘                  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Data Source

Persona narratives are **hardcoded per persona** in the frontend. Each persona ID maps to a story object:

```javascript
const personaStories = {
  priya: {
    narrative: "Priya has created an FD every single year since 2016...",
    hook: "Her CIBIL score? Not Available.",
    punchline: "Traditional credit scoring can't see her. Samridhi can."
  },
  vikram: {
    narrative: "Vikram optimizes everything. 8 different banks, staggered maturities, rate-shopping across institutions...",
    hook: "He's done more financial planning than most loan applicants.",
    punchline: "His FD behavior proves he'd be a model borrower."
  },
  ramesh: {
    narrative: "Ramesh had a good streak going — annual FDs from 2017 to 2019. Then COVID hit his shop...",
    hook: "He broke two FDs to survive. The system would penalize him for it.",
    punchline: "Samridhi sees more — his last 3 deposits show a powerful comeback."
  },
  anita: {
    narrative: "Anita is a consistent saver — she creates FDs regularly. But she breaks them just as often...",
    hook: "She's not irresponsible. She lacks liquidity. She uses her FDs as an emergency fund.",
    punchline: "Give her a small credit line, and she'll stop breaking FDs entirely."
  },
  sunita: {
    narrative: "Sunita retired 5 years ago. She's been saving at SBI — and only SBI — for 15 years...",
    hook: "₹13 lakh in deposits. All at one bank. All in the same tenure.",
    punchline: "Her discipline is flawless. Her diversification could improve — and we'll show her how."
  }
};
```

### Stats Cards
The 4 stats (Deposits, Corpus, Years, Active) are computed client-side from the persona's deposit array. They animate with count-up when the page enters.

### Transition
- Page enters with `fadeSlideUp` animation (opacity 0→1, translateY 30px→0)
- "Analyze" button has a pulse animation
- On click: button changes to a loading state → "Analyzing 9 years of deposit history..." with a spinner → after API response arrives → transition to Score Reveal

### Loading State (Between Profile → Score Reveal)

```
┌─────────────────────────────────────────┐
│                                         │
│         ◌ ─── ◌ ─── ◌ ─── ◌            │
│                                         │
│    Analyzing 9 years of FD history...   │
│                                         │
│    ▸ Evaluating consistency             │
│    ▸ Checking maturity discipline        │
│    ▸ Measuring growth trajectory         │
│    ▸ Assessing diversification           │
│    ▸ Testing tenure intelligence         │
│                                         │
└─────────────────────────────────────────┘
```

This is a **fake loading screen** — the API call takes <100ms, but we deliberately delay for 2.5 seconds and animate each checklist item appearing (500ms apart) to build anticipation. Each step gets a checkmark ✓ animation as it "completes."

**This is the single most impactful UX trick in the entire redesign.** It makes the scoring feel like something is *happening*, not just a function call.

---

## Screen 3: Score Reveal (Complete Redesign)

This is the heart of the experience. Instead of a flat page with everything visible, it's a **scrollable story** with distinct sections that reveal as you scroll.

### Section 3A: The Score Hero

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│     Priya Sharma's FD Credit Score                          │
│                                                             │
│              ┌─────────────────┐                            │
│              │   ╭─────────╮   │                            │
│              │  ╱           ╲  │                            │
│              │ │    689      │  │                            │
│              │  ╲   / 900  ╱  │                            │
│              │   ╰─────────╯   │                            │
│              │     ● GOOD ●    │                            │
│              └─────────────────┘                            │
│                                                             │
│     ┌────────────────┐   ┌─────────────────────┐            │
│     │ CIBIL Equiv:   │   │ Better than 75%     │            │
│     │ 700-749 range  │   │ of savers your age  │            │
│     └────────────────┘   └─────────────────────┘            │
│                                                             │
│              Scroll to see your full breakdown ↓            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Design details:**
- Score dial animates from 0 → 689 over 2 seconds with easing
- The dial's stroke color transitions through bands: red → orange → yellow → green → teal
- Score band badge ("GOOD") pulses gently after the counter completes
- Two context pills appear below with fade-in:
  - CIBIL equivalent: "Equivalent to CIBIL 700-749"
  - Peer percentile: "Better than 75% of savers your age"
- Scroll indicator bounces at the bottom

### Section 3B: Score Radar + Component Cards

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  Score Breakdown                                            │
│                                                             │
│  ┌────────────────────────┐  ┌───────────────────────────┐  │
│  │                        │  │                           │  │
│  │     Radar Chart        │  │  Consistency   ████░ 78   │  │
│  │     (existing)         │  │  30% weight               │  │
│  │                        │  │  ─────────────────────    │  │
│  │                        │  │  FDs/year: 1.2  ██░ 32   │  │
│  │                        │  │  Gap score:     ███ 25   │  │
│  │                        │  │  Streak:        ███ 30   │  │
│  │                        │  │                           │  │
│  └────────────────────────┘  │  Discipline    █████ 95   │  │
│                              │  25% weight               │  │
│                              │  ─────────────────────    │  │
│                              │  Completion:   ████ 40   │  │
│                              │  Withdrawals:  ███░ 28   │  │
│                              │  Hold ratio:   ███░ 27   │  │
│                              │                           │  │
│                              │  [+ 3 more components]    │  │
│                              └───────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Component cards** are individually expandable. Each shows:
- Component name + weight as a tag ("30%")
- Score bar (colored segment out of 100) 
- Sub-metrics listed when expanded, each with its own mini progress bar
- A one-line explanation: "How consistently you create new FDs each year"

The sub-metrics data already exists in the API response (`sub_metrics` array on each component) — we just need to render it.

### Section 3C: Your Pattern (The Insight Moment)

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  Your Savings Pattern                                       │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                                                      │   │
│  │   🏦  LOYAL SINGLE-BANK SAVER                        │   │
│  │                                                      │   │
│  │   Your savings discipline is excellent — 9 years of  │   │
│  │   perfect maturity completion is rare.               │   │
│  │                                                      │   │
│  │   But keeping everything at SBI means you might be   │   │
│  │   missing better rates at HDFC (7.0%), IndusInd      │   │
│  │   (7.5%), or Bandhan (7.5%). You could also be       │   │
│  │   exceeding the DICGC insurance limit of ₹5 lakh    │   │
│  │   per bank.                                          │   │
│  │                                                      │   │
│  │   Spreading across 2-3 banks would boost your score  │   │
│  │   by ~35 points and protect your deposits.           │   │
│  │                                                      │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Design:** This is a large glassmorphism card with a colored left border based on pattern type:
- `DISCIPLINED_OPTIMIZER` → green/teal border
- `LOYAL_SINGLE_BANK` → blue border  
- `LIQUIDITY_GAP_SAVER` → amber border
- `RECOVERING_SAVER` → orange border
- `STANDARD_SAVER` → gray border

The pattern narrative comes from the `insights` array in the API response, but we **enhance it in the frontend** with richer per-pattern text hardcoded in JS (the API's template text is too short for the visual space).

### Section 3D: Deposit Timeline (NEW visualization)

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  Your FD Journey                                            │
│                                                             │
│  ₹2.5L │                                          ●        │
│         │                                    ●  ──┘         │
│  ₹2.0L │                              ●──┘                 │
│         │                         ●──┘                      │
│  ₹1.5L │                    ●──┘                            │
│         │               ●──┘                                │
│  ₹1.0L │          ●──┘                                     │
│         │     ●──┘                                          │
│  ₹0.5L │●──┘                                               │
│         └──────────────────────────────────────────────      │
│          2016  2017  2018  2019  2020  2021  2022  2023  24 │
│                                                             │
│   ● Matured  ◐ Active  ✕ Broken Early                       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Implementation:** Use Chart.js (already loaded) with a **bar chart** showing annual deposit totals, with each bar containing stacked segments:
- Green segments = matured deposits
- Teal segments = active deposits  
- Red/amber segments = withdrawn early

This makes broken FDs visually obvious (Anita's chart will have red segments; Priya's will be all green).

**Data source:** Computed client-side from the persona's deposit array — group by year from `start_date`, sum amounts, color by status.

### Section 3E: Credit Products Unlocked

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  Credit Products Unlocked                                   │
│                                                             │
│  ┌─────────────────┐  ┌──────────────────┐  ┌────────────┐ │
│  │ 💳              │  │ 🏦               │  │ 📱         │ │
│  │ FD-Backed       │  │ FD-Based Loan    │  │ UPI Credit │ │
│  │ Credit Card     │  │                  │  │            │ │
│  │                 │  │ Limit: ₹3.0L     │  │ Limit:     │ │
│  │ Limit: ₹80K    │  │ Rate: 9.5%       │  │ ₹50,000   │ │
│  │ Rate: 18%      │  │                  │  │ Rate: 24%  │ │
│  │                 │  │ ✅ Eligible      │  │            │ │
│  │ ✅ Instant      │  │ Quick disbursal  │  │ ✅ Ready   │ │
│  │   Approval     │  │                  │  │            │ │
│  └─────────────────┘  └──────────────────┘  └────────────┘ │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ 🔒 Personal Loan — ₹5L limit                        │   │
│  │    Requires score ≥ 600  ✅ You qualify (689)        │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Design:** Each product is a premium glassmorphism card with:
- Icon (emoji for now)
- Product name
- Key terms (limit, rate)
- Eligibility badge: green ✅ if eligible, amber ⚠️ if close, red 🔒 if locked
- The locked products show "Score 750+ to unlock" — creates a gamification hook

### Section 3F: Improvement Path (with Score Projection)

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  Your Improvement Path                                      │
│                                                             │
│  Current: 689 (Good)  →  Potential: 739 (Very Good)         │
│  ████████████████████████████░░░░░░░░                       │
│                         ↑ you are here  ↑ potential         │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  1. Open an FD at a second bank                      │   │
│  │     Difficulty: Easy  •  Impact: +35 pts             │   │
│  │     ─────────────────────────────── ████████████      │   │
│  │                                                      │   │
│  │  2. Maintain regular annual deposits                 │   │
│  │     Difficulty: Easy  •  Impact: +25 pts             │   │
│  │     ─────────────────────────────── █████████        │   │
│  │                                                      │   │
│  │  3. Try a shorter tenure (6-month FD)                │   │
│  │     Difficulty: Medium  •  Impact: +15 pts           │   │
│  │     ─────────────────────────────── ██████           │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Key design:**
- **Score projection bar**: Shows current score position AND potential score position on a 300-900 scale bar, with the gap highlighted in gradient
- Each improvement has a difficulty badge (color-coded: green Easy, yellow Medium, red Hard) and an impact bar
- The projected score is computed backend: `current_score + sum(top_3_improvements.points_delta)`

### Section 3G: The Flywheel (Blostem Value Prop)

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│  The Samridhi Flywheel                                      │
│                                                             │
│      Save in FDs  ──→  Score improves  ──→  Better credit   │
│           ↑                                       │         │
│           └──── Use credit wisely ←───────────────┘         │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  "Every FD you create builds your credit identity.   │   │
│  │   Every maturity you honor proves your discipline.   │   │
│  │   Samridhi turns your savings history into your      │   │
│  │   credit future."                                    │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                             │
│          [← Try Another Persona]  [Your Own Score →]        │
│                                                             │
│     Built for the Blostem AI Builder Hackathon 2026         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

This is a simple CSS animation showing the 4-step flywheel cycle. It can be done with CSS `@keyframes`:
- 4 nodes in a circle/diamond layout
- An animated path/glow that travels around the cycle
- A quote block below with the value proposition

---

## Screen 4: Form Wizard (Complete Redesign of Custom Input)

### Current Problem
The custom input asks users to type raw JSON. Nobody will do this.

### New Design
A multi-step wizard:

**Step 1: About You**
```
┌──────────────────────────────────────┐
│  Step 1 of 3 — About You            │
│  ━━━━━━━━━━━━━░░░░░░░░░░            │
│                                      │
│  Name:  [___________________]        │
│  Age:   [___]                        │
│  City:  [___________________]        │
│                                      │
│              [Next →]                │
└──────────────────────────────────────┘
```

**Step 2: Add Your Deposits** (repeatable)
```
┌──────────────────────────────────────┐
│  Step 2 of 3 — Your Deposits         │
│  ━━━━━━━━━━━━━━━━━━━░░░░            │
│                                      │
│  Deposit #1                          │
│  Type:    [FD ▾]    [RD ▾]           │
│  Bank:    [SBI ▾]  (dropdown)        │
│  Amount:  [₹ ________]              │
│  Tenure:  [12 months ▾]             │
│  Rate:    [6.5 %]                   │
│  Started: [2024-01-01]              │
│  Status:  [Active ▾] [Matured ▾]    │
│                                      │
│  ┌─────────────────────────────┐     │
│  │  ✓ ₹1,00,000 FD at SBI     │     │
│  │    12 months, 6.5%, Active  │     │
│  └─────────────────────────────┘     │
│                                      │
│  [+ Add Another Deposit]            │
│                                      │
│       [← Back]  [Next →]            │
└──────────────────────────────────────┘
```

**Key UX decisions:**
- Bank is a **dropdown** with the 14 banks already in the `bankRates` map
- Tenure is a dropdown: 3mo, 6mo, 12mo, 24mo, 36mo, 60mo
- Status dropdown: Active / Matured / Withdrawn Early
- If "Withdrawn Early" → show date picker for withdrawn date
- Each added deposit appears as a summary card below the form
- Minimum 1 deposit required to proceed

**Step 3: Review & Calculate**
```
┌──────────────────────────────────────┐
│  Step 3 of 3 — Review                │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━       │
│                                      │
│  You: Test User, 30, Delhi           │
│  Deposits: 3                         │
│                                      │
│  1. ₹1,00,000 FD at SBI (Active)    │
│  2. ₹50,000 FD at HDFC (Matured)    │
│  3. ₹5,000/mo RD at ICICI (Active)  │
│                                      │
│     [← Edit]  [✨ Calculate Score]   │
└──────────────────────────────────────┘
```

After clicking "Calculate Score" → same loading ceremony → score reveal.

---

## CSS Design System Enhancements

### New Design Tokens

```css
:root {
  /* Existing tokens preserved */
  
  /* NEW: Score band colors */
  --score-excellent: #00d4aa;
  --score-good: #4ecdc4;
  --score-fair: #ffe66d;
  --score-needs-work: #ff8a5c;
  --score-poor: #ff6b6b;
  
  /* NEW: Pattern colors */
  --pattern-optimizer: #00d4aa;
  --pattern-loyal: #4ea8de;
  --pattern-liquidity: #f59e0b;
  --pattern-recovering: #f97316;
  --pattern-standard: #6b7280;
  
  /* NEW: Difficulty badges */
  --diff-easy: #00d4aa;
  --diff-medium: #f59e0b;
  --diff-hard: #ef4444;
  
  /* NEW: Spacing scale */
  --section-gap: 64px;
  --card-gap: 20px;
  
  /* NEW: Shadows */
  --shadow-card: 0 4px 24px rgba(0, 0, 0, 0.2);
  --shadow-hover: 0 8px 40px rgba(0, 212, 170, 0.15);
  --shadow-glow: 0 0 30px rgba(0, 212, 170, 0.3);
}
```

### New Animations

```css
/* Staggered fade-up for lists */
@keyframes fadeSlideUp {
  from { opacity: 0; transform: translateY(24px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Score band pulse */
@keyframes bandPulse {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.05); opacity: 0.9; }
}

/* Loading checkmark */
@keyframes checkIn {
  from { opacity: 0; transform: scale(0.5) rotate(-10deg); }
  to { opacity: 1; transform: scale(1) rotate(0); }
}

/* Scroll indicator bounce */
@keyframes scrollBounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(8px); }
}

/* Flywheel rotation glow */
@keyframes flywheelGlow {
  0% { box-shadow: 0 0 20px var(--accent-primary); }
  25% { box-shadow: 0 0 5px transparent; }
  100% { box-shadow: 0 0 20px var(--accent-primary); }
}

/* Progress bar fill */
@keyframes fillBar {
  from { width: 0; }
  to { width: var(--bar-width); }
}
```

### Scroll-triggered animations
Use `IntersectionObserver` to trigger animations when sections scroll into view:

```javascript
const observer = new IntersectionObserver((entries) => {
  entries.forEach(entry => {
    if (entry.isIntersecting) {
      entry.target.classList.add('animate-in');
      observer.unobserve(entry.target);
    }
  });
}, { threshold: 0.15 });

document.querySelectorAll('.scroll-reveal').forEach(el => {
  observer.observe(el);
});
```

Each section of the score reveal page has `class="scroll-reveal"` and only animates when scrolled into viewport.

---

## Rebranding: "FD Credit Score" → "Samridhi"

The project repo is already named `samri-dhi`. Let's use this as the product name:

- **Samridhi** (समृद्धि) = Prosperity in Hindi/Sanskrit
- Logo text: "◈ Samridhi" 
- Tagline: "Your savings. Your credit identity."
- This gives it a product identity rather than a feature description

Update in: header, page title, all copy, footer.

---

## File Changes Summary

### Frontend (all in `web/`)

| File | Action | What Changes |
|------|--------|-------------|
| `index.html` | **REWRITE** | Complete restructure with all 4 screens, new sections, semantic HTML |
| `css/style.css` | **REWRITE** | Full design system with new tokens, components, animations, responsive |
| `js/app.js` | **REWRITE** | New view management, persona stories, form wizard, loading ceremony, scroll animations |
| `js/charts.js` | **MODIFY** | Add deposit timeline chart alongside existing radar chart |

### Backend (in `internal/`)

| File | Action | What Changes |
|------|--------|-------------|
| `handlers/api.go` | **MODIFY** | Wrap persona score response with persona metadata + computed stats |
| `scoring/engine.go` | **MODIFY** | Add `projected_score` and `cibil_equivalent` to ScoreResult |
| `scoring/models.go` | **MODIFY** | Add `ProjectedScore`, `CIBILEquivalent` fields to ScoreResult |

---

## Execution Timeline

| Phase | Hours | Work |
|-------|-------|------|
| **1. Backend tweak** | 0.5h | Add projected_score, cibil_equivalent, persona stats to API |
| **2. HTML structure** | 1.5h | Rewrite index.html with all 4 screens and sections |
| **3. CSS design system** | 2h | Complete style.css rewrite — tokens, components, animations, responsive |
| **4. Landing page JS** | 1h | Stats counter animation, persona cards with taglines, routing |
| **5. Profile page JS** | 1h | Persona stories, stat cards, loading ceremony |
| **6. Score reveal JS** | 2h | Score animation, component cards, pattern narrative, products, improvements |
| **7. Deposit timeline** | 1h | Chart.js bar chart for deposit history |  
| **8. Form wizard** | 1.5h | 3-step wizard with deposit form, validation, review |
| **9. Polish** | 1h | Scroll animations, transitions, edge cases, mobile |
| **10. Test + Deploy** | 1h | Build, test all personas, deploy |
| **Total** | **~13h** | |

---

## Verification Plan

### Manual Testing
- Click every persona → verify profile page shows correct story → score reveal loads
- Verify all 5 persona scores match expected ranges from the deep-dive document
- Test form wizard: add 1 deposit, add 3 deposits, test validation
- Test mobile viewport (375px, 768px)
- Test scroll animations fire correctly
- Verify deposit timeline shows correct data per persona

### Browser Testing
- Use browser tool to navigate full flow and capture recording
- Click Priya → profile → analyze → score reveal → scroll through all sections
- Click Anita → verify "Liquidity Gap Saver" pattern shows with amber styling
- Test custom input form wizard end-to-end

### Visual Verification
- Screenshot each major screen section
- Verify dark theme consistency
- Verify animations are smooth (no jank)
- Verify glassmorphism effects render correctly

---

## Open Questions

> [!IMPORTANT]
> **Go installation**: Go is not installed in this environment. Options:
> 1. Install Go and rebuild the binary
> 2. Use the existing pre-built `fd-credit-score` binary (8.7MB, already in repo root)
> 3. Skip backend changes and hardcode the new fields entirely in the frontend
>
> **Recommendation**: Option 2 first (test existing binary), then Option 1 if backend changes are needed. The backend changes are small enough that we could also fake them entirely in the frontend JS (compute persona stats and projected score client-side).

> [!IMPORTANT]
> **Deployment**: Where will this be deployed? The Go binary serves everything, so any platform that can run a Linux binary works (Render, Railway, Fly.io, or even a free Oracle Cloud VM).
