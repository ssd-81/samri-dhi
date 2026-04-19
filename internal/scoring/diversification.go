package scoring

import (
	"math"
	"sort"
)

func calculateDiversification(history DepositHistory) ComponentScore {
	if len(history.Deposits) == 0 {
		return ComponentScore{Name: "Diversification", Score: 0, MaxScore: 100, Weight: 0.15}
	}

	bankCount := calculateUniqueBanks(history.Deposits)
	bankScore := 0
	if bankCount >= 5 {
		bankScore = 35
	} else if bankCount >= 3 {
		bankScore = 25
	} else if bankCount == 2 {
		bankScore = 15
	} else {
		bankScore = 5
	}

	tenureSpreadScore := calculateTenureSpread(history.Deposits)

	productMixScore := calculateProductMix(history.Deposits)

	totalScore := bankScore + tenureSpreadScore + productMixScore
	if totalScore > 100 {
		totalScore = 100
	}

	subMetrics := []SubMetric{
		{Name: "Bank Count", Value: float64(bankCount), Score: bankScore, Max: 35},
		{Name: "Tenure Spread", Value: 0, Score: tenureSpreadScore, Max: 35},
		{Name: "Product Mix", Value: 0, Score: productMixScore, Max: 30},
	}

	return ComponentScore{
		Name:       "Diversification",
		Score:      totalScore,
		MaxScore:   100,
		Weight:     0.15,
		Weighted:   float64(totalScore) * 0.15,
		SubMetrics: subMetrics,
	}
}

func calculateUniqueBanks(deposits []Deposit) int {
	banks := make(map[string]bool)
	for _, d := range deposits {
		if d.Bank != "" {
			banks[d.Bank] = true
		}
	}
	return len(banks)
}

func calculateTenureSpread(deposits []Deposit) int {
	var tenures []int
	for _, d := range deposits {
		if d.TenureMonths > 0 {
			tenures = append(tenures, d.TenureMonths)
		}
	}

	if len(tenures) < 2 {
		return 10
	}

	mean := 0.0
	for _, t := range tenures {
		mean += float64(t)
	}
	mean /= float64(len(tenures))

	variance := 0.0
	for _, t := range tenures {
		diff := float64(t) - mean
		variance += diff * diff
	}
	variance /= float64(len(tenures))

	stdDev := math.Sqrt(variance)
	coefficientOfVariation := stdDev / mean

	if coefficientOfVariation >= 0.5 {
		return 35
	} else if coefficientOfVariation >= 0.3 {
		return 25
	} else if coefficientOfVariation >= 0.15 {
		return 15
	}
	return 5
}

func calculateProductMix(deposits []Deposit) int {
	hasFD := false
	hasRD := false
	hasTaxSaver := false

	for _, d := range deposits {
		if d.Type == "FD" {
			hasFD = true
			if d.TenureMonths >= 60 {
				hasTaxSaver = true
			}
		}
		if d.Type == "RD" {
			hasRD = true
		}
	}

	score := 0
	if hasFD {
		score += 10
	}
	if hasRD {
		score += 10
	}
	if hasTaxSaver {
		score += 10
	}

	uniqueTypes := 0
	if hasFD {
		uniqueTypes++
	}
	if hasRD {
		uniqueTypes++
	}

	if uniqueTypes >= 2 {
		score += 10
	}

	return score
}

var _ = sort.Ints