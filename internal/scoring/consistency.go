package scoring

import (
	"sort"
	"strconv"
	"time"
)

func calculateConsistency(history DepositHistory) ComponentScore {
	if len(history.Deposits) == 0 {
		return ComponentScore{Name: "Consistency", Score: 0, MaxScore: 100, Weight: 0.30}
	}

	deposits := history.Deposits
	fdCount := 0
	for _, d := range deposits {
		if d.Type == "FD" {
			fdCount++
		}
	}

	years := calculateActiveYears(deposits)
	if years == 0 {
		years = 1
	}

	depositFreqScore := normalize(fdCount, years, 0, 40)
	subMetrics := []SubMetric{
		{Name: "FDs per year", Value: float64(fdCount) / float64(years), Score: depositFreqScore, Max: 40},
	}

	gapScore := calculateGapScore(deposits)
	subMetrics = append(subMetrics, SubMetric{Name: "Gap Analysis", Value: 0, Score: gapScore, Max: 30})

	streakScore := calculateStreakScore(deposits)
	subMetrics = append(subMetrics, SubMetric{Name: "Streak", Value: 0, Score: streakScore, Max: 30})

	totalScore := depositFreqScore + gapScore + streakScore

	return ComponentScore{
		Name:       "Consistency",
		Score:      totalScore,
		MaxScore:   100,
		Weight:     0.30,
		Weighted:   float64(totalScore) * 0.30,
		SubMetrics: subMetrics,
	}
}

func calculateActiveYears(deposits []Deposit) int {
	if len(deposits) == 0 {
		return 0
	}
	dates := extractDates(deposits)
	if len(dates) == 0 {
		return 0
	}
	minYear, maxYear := dates[0].Year(), dates[len(dates)-1].Year()
	if maxYear-minYear < 1 {
		return 1
	}
	return maxYear - minYear + 1
}

func extractDates(deposits []Deposit) []time.Time {
	var dates []time.Time
	for _, d := range deposits {
		if t, err := time.Parse("2006-01-02", d.StartDate); err == nil {
			dates = append(dates, t)
		}
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})
	return dates
}

func normalize(value, total, min, max int) int {
	if total == 0 {
		return min
	}
	ratio := float64(value) / float64(total)
	if ratio > 1 {
		ratio = 1
	}
	return min + int(ratio*float64(max-min))
}

func calculateGapScore(deposits []Deposit) int {
	dates := extractDates(deposits)
	if len(dates) < 2 {
		return 30
	}

	var gaps []int
	for i := 1; i < len(dates); i++ {
		months := int(dates[i].Sub(dates[i-1]).Hours() / (24 * 30))
		if months > 0 && months < 24 {
			gaps = append(gaps, months)
		}
	}

	if len(gaps) == 0 {
		return 15
	}

	avg := 0
	for _, g := range gaps {
		avg += g
	}
	avg /= len(gaps)

	if avg <= 2 {
		return 30
	} else if avg <= 4 {
		return 20
	} else if avg <= 6 {
		return 10
	}
	return 5
}

func calculateStreakScore(deposits []Deposit) int {
	years := make(map[int]int)
	for _, d := range deposits {
		if t, err := time.Parse("2006-01-02", d.StartDate); err == nil {
			years[t.Year()]++
		}
	}

	streak := 0
	currentStreak := 0
	prevYear := 0
	for y := 2010; y <= 2026; y++ {
		if years[y] > 0 {
			if y == prevYear+1 || prevYear == 0 {
				currentStreak++
			} else {
				if currentStreak > streak {
					streak = currentStreak
				}
				currentStreak = 1
			}
			prevYear = y
		} else {
			if currentStreak > streak {
				streak = currentStreak
			}
			currentStreak = 0
		}
	}
	if currentStreak > streak {
		streak = currentStreak
	}

	if streak >= 7 {
		return 30
	} else if streak >= 5 {
		return 25
	} else if streak >= 3 {
		return 20
	} else if streak >= 1 {
		return 10
	}
	return 0
}

func calculateYearFromDate(dateStr string) int {
	if len(dateStr) >= 4 {
		if y, err := strconv.Atoi(dateStr[:4]); err == nil {
			return y
		}
	}
	return 0
}