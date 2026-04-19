package scoring

import (
	"sort"
)

func calculateGrowth(history DepositHistory) ComponentScore {
	if len(history.Deposits) == 0 {
		return ComponentScore{Name: "Growth Trajectory", Score: 0, MaxScore: 100, Weight: 0.20}
	}

	yearTotals := make(map[int]float64)
	for _, d := range history.Deposits {
		year := calculateYearFromDate(d.StartDate)
		if year > 0 {
			yearTotals[year] += d.Amount
		}
	}

	var years []int
	for y := range yearTotals {
		years = append(years, y)
	}
	sort.Ints(years)

	growthScore := 0
	recoveryScore := 0

	if len(years) >= 2 {
		firstHalf := 0.0
		secondHalf := 0.0
		mid := len(years) / 2

		for i := 0; i < mid; i++ {
			firstHalf += yearTotals[years[i]]
		}
		for i := mid; i < len(years); i++ {
			secondHalf += yearTotals[years[i]]
		}

		if firstHalf > 0 {
			growth := ((secondHalf - firstHalf) / firstHalf) * 100
			if growth >= 30 {
				growthScore = 40
			} else if growth >= 15 {
				growthScore = 30
			} else if growth >= 5 {
				growthScore = 20
			} else if growth >= 0 {
				growthScore = 10
			} else {
				growthScore = 0
			}
		}

		hasDip := false
		hasRecovery := false
		for i := 1; i < len(years); i++ {
			if yearTotals[years[i]] < yearTotals[years[i-1]]*0.7 {
				hasDip = true
			}
			if hasDip && yearTotals[years[i]] > yearTotals[years[i-1]]*1.2 {
				hasRecovery = true
			}
		}
		if hasDip && hasRecovery {
			recoveryScore = 20
		} else if hasRecovery {
			recoveryScore = 10
		}
	}

	corpusScore := calculateCorpusScore(history)

	totalScore := growthScore + recoveryScore + corpusScore
	if totalScore > 100 {
		totalScore = 100
	}

	subMetrics := []SubMetric{
		{Name: "YoY Growth", Value: 0, Score: growthScore, Max: 40},
		{Name: "Recovery Pattern", Value: 0, Score: recoveryScore, Max: 20},
		{Name: "Corpus Size", Value: 0, Score: corpusScore, Max: 40},
	}

	return ComponentScore{
		Name:       "Growth Trajectory",
		Score:      totalScore,
		MaxScore:   100,
		Weight:     0.20,
		Weighted:   float64(totalScore) * 0.20,
		SubMetrics: subMetrics,
	}
}

func calculateCorpusScore(history DepositHistory) int {
	var activeTotal float64
	for _, d := range history.Deposits {
		if d.Status == "active" {
			activeTotal += d.Amount
		}
	}

	if activeTotal >= 3000000 {
		return 40
	} else if activeTotal >= 1500000 {
		return 30
	} else if activeTotal >= 500000 {
		return 20
	} else if activeTotal >= 200000 {
		return 10
	}
	return 5
}