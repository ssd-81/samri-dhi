package scoring

import "time"

func calculateDiscipline(history DepositHistory) ComponentScore {
	if len(history.Deposits) == 0 {
		return ComponentScore{Name: "Maturity Discipline", Score: 0, MaxScore: 100, Weight: 0.25}
	}

	deposits := history.Deposits
	var closedFDs []Deposit
	var maturedFDs int
	var withdrawnEarly int
	var totalHoldMonths int
	var holdCount int

	for _, d := range deposits {
		if d.Type == "FD" && d.Status != "active" {
			closedFDs = append(closedFDs, d)
			if d.Status == "matured" || (d.WithdrawnDate == nil && isMatured(d)) {
				maturedFDs++
			} else if d.WithdrawnDate != nil {
				withdrawnEarly++
			}

			if d.Status == "matured" || d.WithdrawnDate != nil {
				start, _ := time.Parse("2006-01-02", d.StartDate)
				var end time.Time
				if d.WithdrawnDate != nil {
					end, _ = time.Parse("2006-01-02", *d.WithdrawnDate)
				} else {
					end, _ = time.Parse("2006-01-02", d.MaturityDate)
				}
				months := int(end.Sub(start).Hours() / (24 * 30))
				if months > 0 {
					holdRatio := float64(months) / float64(d.TenureMonths)
					if holdRatio > 1 {
						holdRatio = 1
					}
					totalHoldMonths += int(holdRatio * 100)
					holdCount++
				}
			}
		}
	}

	completionRate := 0.0
	if len(closedFDs) > 0 {
		completionRate = float64(maturedFDs) / float64(len(closedFDs))
	}
	completionScore := int(completionRate * 40)

	withdrawalPenalty := 0
	if withdrawnEarly > 0 {
		if withdrawnEarly >= 3 {
			withdrawalPenalty = 30
		} else if withdrawnEarly == 2 {
			withdrawalPenalty = 20
		} else {
			withdrawalPenalty = 10
		}
	}

	holdScore := 0
	if holdCount > 0 {
		holdScore = totalHoldMonths / holdCount
	} else {
		holdScore = 30
	}

	rdScore := calculateRDDiscipline(history)

	totalScore := completionScore + (30 - withdrawalPenalty) + holdScore + rdScore
	if totalScore > 100 {
		totalScore = 100
	}

	subMetrics := []SubMetric{
		{Name: "Completion Rate", Value: completionRate * 100, Score: completionScore, Max: 40},
		{Name: "Early Withdrawals", Value: float64(withdrawnEarly), Score: 30 - withdrawalPenalty, Max: 30},
		{Name: "Hold Ratio", Value: 0, Score: holdScore, Max: 30},
		{Name: "RD Discipline", Value: 0, Score: rdScore, Max: 10},
	}

	return ComponentScore{
		Name:       "Maturity Discipline",
		Score:      totalScore,
		MaxScore:   100,
		Weight:     0.25,
		Weighted:   float64(totalScore) * 0.25,
		SubMetrics: subMetrics,
	}
}

func calculateRDDiscipline(history DepositHistory) int {
	var rdDeposits []Deposit
	for _, d := range history.Deposits {
		if d.Type == "RD" {
			rdDeposits = append(rdDeposits, d)
		}
	}

	if len(rdDeposits) == 0 {
		return 0
	}

	completed := 0
	for _, d := range rdDeposits {
		if d.RDPaidMonths >= d.TenureMonths-1 {
			completed++
		} else if d.RDMissedMonths > 0 {
			missRate := float64(d.RDMissedMonths) / float64(d.RDPaidMonths+d.RDMissedMonths)
			if missRate < 0.1 {
				completed++
			}
		}
	}

	return (completed * 100) / (len(rdDeposits) * 10)
}

func isMatured(d Deposit) bool {
	if d.Status == "matured" {
		return true
	}
	if d.MaturityDate == "" {
		return false
	}
	maturity, err := time.Parse("2006-01-02", d.MaturityDate)
	if err != nil {
		return false
	}
	return maturity.Before(time.Now())
}

var _ = time.Now