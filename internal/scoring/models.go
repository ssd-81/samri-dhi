package scoring

type DepositHistory struct {
	UserID   string     `json:"user_id"`
	Name     string     `json:"name"`
	Age      int        `json:"age"`
	City     string     `json:"city"`
	Deposits []Deposit  `json:"deposits"`
}

type Deposit struct {
	ID             string   `json:"id"`
	Type           string   `json:"type"`
	Bank           string   `json:"bank"`
	Amount         float64  `json:"amount"`
	TenureMonths   int      `json:"tenure_months"`
	InterestRate   float64  `json:"interest_rate"`
	StartDate      string   `json:"start_date"`
	MaturityDate   string   `json:"maturity_date"`
	WithdrawnDate  *string  `json:"withdrawn_date"`
	Status         string   `json:"status"`
	RDInstallment  float64  `json:"rd_installment,omitempty"`
	RDPaidMonths   int      `json:"rd_paid_months,omitempty"`
	RDMissedMonths int      `json:"rd_missed_months,omitempty"`
}

type ScoreResult struct {
	TotalScore      int              `json:"total_score"`
	ScoreBand       string           `json:"score_band"`
	ProjectedScore  int              `json:"projected_score"`
	CIBILEquivalent string           `json:"cibil_equivalent"`
	Components      []ComponentScore `json:"components"`
	CreditProducts  []CreditProduct  `json:"credit_products"`
	Improvements    []Improvement    `json:"improvements"`
	Insights        []string         `json:"insights"`
	PatternDetected string           `json:"pattern_detected"`
	PeerPercentile  int              `json:"peer_percentile"`
}

type PersonaSummary struct {
	Name           string `json:"name"`
	Age            int    `json:"age"`
	Occupation     string `json:"occupation"`
	City           string `json:"city"`
	DepositCount   int    `json:"deposit_count"`
	TotalCorpus    float64 `json:"total_corpus"`
	YearsActive    int    `json:"years_active"`
	ActiveDeposits int    `json:"active_deposits"`
}

type ComponentScore struct {
	Name       string      `json:"name"`
	Score      int         `json:"score"`
	MaxScore   int         `json:"max_score"`
	Weight     float64    `json:"weight"`
	Weighted   float64    `json:"weighted"`
	SubMetrics []SubMetric `json:"sub_metrics"`
}

type SubMetric struct {
	Name   string  `json:"name"`
	Value  float64 `json:"value"`
	Score  int     `json:"score"`
	Max    int     `json:"max"`
}

type CreditProduct struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Limit       float64 `json:"limit"`
	Interest    float64 `json:"interest"`
	Eligibility string  `json:"eligibility"`
}

type Improvement struct {
	Action      string `json:"action"`
	PointsDelta int    `json:"points_delta"`
	Difficulty  string `json:"difficulty"`
}