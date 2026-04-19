package scoring

import (
	"sort"
	"time"
)

func calculateIntelligence(history DepositHistory) ComponentScore {
	if len(history.Deposits) == 0 {
		return ComponentScore{Name: "Tenure Intelligence", Score: 0, MaxScore: 100, Weight: 0.10}
	}

	rateScore := calculateRateOptimization(history.Deposits)
	ladderScore := calculateLadderDetection(history.Deposits)
	taxScore := calculateTaxSaverAwareness(history.Deposits)

	totalScore := rateScore + ladderScore + taxScore
	if totalScore > 100 {
		totalScore = 100
	}

	subMetrics := []SubMetric{
		{Name: "Rate Optimization", Value: 0, Score: rateScore, Max: 40},
		{Name: "Ladder Detection", Value: 0, Score: ladderScore, Max: 30},
		{Name: "Tax-Saver Awareness", Value: 0, Score: taxScore, Max: 30},
	}

	return ComponentScore{
		Name:       "Tenure Intelligence",
		Score:      totalScore,
		MaxScore:   100,
		Weight:     0.10,
		Weighted:   float64(totalScore) * 0.10,
		SubMetrics: subMetrics,
	}
}

var bankRates = map[string]float64{
	"SBI":        6.8,
	"HDFC":       7.0,
	"ICICI":      7.0,
	"Axis":       7.2,
	"BoB":        6.5,
	"Canara":     6.5,
	"PNB":        6.5,
	"IDBI":       6.5,
	"Indian":     6.8,
	"Kotak":      7.0,
	"IndusInd":   7.5,
	"Yes":        7.0,
	"RBL":        7.5,
	"Bandhan":    7.5,
}

func calculateRateOptimization(deposits []Deposit) int {
	highRateCount := 0
	totalFDs := 0

	for _, d := range deposits {
		if d.Type == "FD" {
			totalFDs++
			if d.InterestRate >= 7.0 {
				highRateCount++
			}
		}
	}

	if totalFDs == 0 {
		return 20
	}

	ratio := float64(highRateCount) / float64(totalFDs)

	if ratio >= 0.8 {
		return 40
	} else if ratio >= 0.5 {
		return 30
	} else if ratio >= 0.3 {
		return 20
	}
	return 10
}

func calculateLadderDetection(deposits []Deposit) int {
	activeFDs := make([]time.Time, 0)
	for _, d := range deposits {
		if d.Type == "FD" && d.Status == "active" {
			if t, err := time.Parse("2006-01-02", d.MaturityDate); err == nil {
				activeFDs = append(activeFDs, t)
			}
		}
	}

	if len(activeFDs) < 2 {
		return 10
	}

	sort.Slice(activeFDs, func(i, j int) bool {
		return activeFDs[i].Before(activeFDs[j])
	})

	gaps := make([]int, 0)
	for i := 1; i < len(activeFDs); i++ {
		months := int(activeFDs[i].Sub(activeFDs[i-1]).Hours() / (24 * 30))
		if months > 0 && months < 12 {
			gaps = append(gaps, months)
		}
	}

	if len(gaps) >= len(activeFDs)/2 {
		return 30
	} else if len(gaps) >= 1 {
		return 20
	}
	return 10
}

func calculateTaxSaverAwareness(deposits []Deposit) int {
	taxSaverCount := 0
	totalLongTerm := 0

	for _, d := range deposits {
		if d.Type == "FD" && d.TenureMonths >= 60 {
			totalLongTerm++
			if d.TenureMonths >= 60 {
				taxSaverCount++
			}
		}
	}

	if totalLongTerm == 0 {
		return 5
	}

	ratio := float64(taxSaverCount) / float64(totalLongTerm)
	if ratio >= 0.5 {
		return 30
	} else if ratio >= 0.3 {
		return 20
	} else if ratio >= 0.1 {
		return 10
	}
	return 5
}

var _ = time.Now