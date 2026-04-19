package scoring

import (
	"sort"
)

func CalculateScore(history DepositHistory) ScoreResult {
	components := []ComponentScore{
		calculateConsistency(history),
		calculateDiscipline(history),
		calculateGrowth(history),
		calculateDiversification(history),
		calculateIntelligence(history),
	}

	var rawScore float64
	for _, c := range components {
		rawScore += c.Weighted
	}

	mappedScore := mapTo900(rawScore)
	pattern := detectPattern(components, history)
	products := recommendProducts(mappedScore, history, pattern)
	improvements := suggestImprovements(components, history)
	insights := generateInsights(components, pattern)
	percentile := calculatePercentile(mappedScore)

	return ScoreResult{
		TotalScore:      mappedScore,
		ScoreBand:       getScoreBand(mappedScore),
		Components:      components,
		CreditProducts:  products,
		Improvements:    improvements,
		Insights:        insights,
		PatternDetected: pattern,
		PeerPercentile:  percentile,
	}
}

func mapTo900(rawScore float64) int {
	score := int(rawScore * 9)
	if score > 900 {
		score = 900
	}
	if score < 300 {
		score = 300
	}
	return score
}

func getScoreBand(score int) string {
	switch {
	case score >= 750:
		return "Excellent"
	case score >= 650:
		return "Good"
	case score >= 550:
		return "Fair"
	case score >= 450:
		return "Needs Improvement"
	default:
		return "Poor"
	}
}

func detectPattern(components []ComponentScore, history DepositHistory) string {
	consistency := findComponent("Consistency", components)
	discipline := findComponent("Maturity Discipline", components)
	diversification := findComponent("Diversification", components)
	growth := findComponent("Growth Trajectory", components)

	if consistency.Score > 65 && discipline.Score < 55 {
		return "LIQUIDITY_GAP_SAVER"
	}

	if avgScore(components) > 70 && diversification.Score < 35 {
		return "LOYAL_SINGLE_BANK"
	}

	if consistency.Score < 40 && growth.Score > 50 {
		return "RECOVERING_SAVER"
	}

	if minScore(components) > 70 {
		return "DISCIPLINED_OPTIMIZER"
	}

	if consistency.Score > 60 && discipline.Score > 70 && diversification.Score < 30 {
		return "CAUTIOUS_DIVERSIFIER"
	}

	return "STANDARD_SAVER"
}

func findComponent(name string, components []ComponentScore) ComponentScore {
	for _, c := range components {
		if c.Name == name {
			return c
		}
	}
	return ComponentScore{}
}

func avgScore(components []ComponentScore) float64 {
	if len(components) == 0 {
		return 0
	}
	sum := 0
	for _, c := range components {
		sum += c.Score
	}
	return float64(sum) / float64(len(components))
}

func minScore(components []ComponentScore) int {
	if len(components) == 0 {
		return 0
	}
	min := components[0].Score
	for _, c := range components {
		if c.Score < min {
			min = c.Score
		}
	}
	return min
}

func calculatePercentile(score int) int {
	switch {
	case score >= 750:
		return 90
	case score >= 650:
		return 75
	case score >= 550:
		return 50
	case score >= 450:
		return 25
	default:
		return 10
	}
}

func recommendProducts(score int, history DepositHistory, pattern string) []CreditProduct {
	var products []CreditProduct
	var activeTotal float64
	for _, d := range history.Deposits {
		if d.Status == "active" {
			activeTotal += d.Amount
		}
	}

	if activeTotal >= 500000 {
		products = append(products, CreditProduct{
			Name:        "FD-Backed Credit Card",
			Type:        "Secured",
			Limit:       activeTotal * 0.2,
			Interest:    18,
			Eligibility: "Instant approval",
		})
	}

	if activeTotal >= 100000 {
		products = append(products, CreditProduct{
			Name:        "FD-Based Loan",
			Type:        "Secured",
			Limit:       activeTotal * 0.75,
			Interest:    9.5,
			Eligibility: "Low interest, quick disbursal",
		})
	}

	products = append(products, CreditProduct{
		Name:        "UPI Credit",
		Type:        "Credit Line",
		Limit:       50000,
		Interest:    24,
		Eligibility: "Based on FD holdings",
	})

	if score >= 600 {
		products = append(products, CreditProduct{
			Name:        "Personal Loan",
			Type:        "Unsecured",
			Limit:       500000,
			Interest:    12,
			Eligibility: "Requires score > 600",
		})
	}

	return products
}

func suggestImprovements(components []ComponentScore, history DepositHistory) []Improvement {
	var improvements []Improvement

	diversification := findComponent("Diversification", components)
	if diversification.Score < 35 {
		improvements = append(improvements, Improvement{
			Action:      "Open FD at a second bank",
			PointsDelta: 35,
			Difficulty:  "Easy",
		})
	}

	discipline := findComponent("Maturity Discipline", components)
	if discipline.Score < 60 {
		improvements = append(improvements, Improvement{
			Action:      "Avoid premature withdrawals",
			PointsDelta: 20,
			Difficulty:  "Medium",
		})
	}

	intelligence := findComponent("Tenure Intelligence", components)
	if intelligence.Score < 40 {
		improvements = append(improvements, Improvement{
			Action:      "Create an FD ladder (stagger maturities)",
			PointsDelta: 15,
			Difficulty:  "Medium",
		})
	}

	consistency := findComponent("Consistency", components)
	if consistency.Score < 50 {
		improvements = append(improvements, Improvement{
			Action:      "Maintain regular annual deposits",
			PointsDelta: 25,
			Difficulty:  "Easy",
		})
	}

	growth := findComponent("Growth Trajectory", components)
	if growth.Score < 40 {
		improvements = append(improvements, Improvement{
			Action:      "Increase deposit amounts gradually",
			PointsDelta: 15,
			Difficulty:  "Hard",
		})
	}

	bankCount := 0
	banks := make(map[string]bool)
	for _, d := range history.Deposits {
		if d.Bank != "" {
			banks[d.Bank] = true
		}
	}
	bankCount = len(banks)

	if bankCount == 1 {
		improvements = append(improvements, Improvement{
			Action:      "Diversify across banks for better rates",
			PointsDelta: 20,
			Difficulty:  "Easy",
		})
	}

	sort.Slice(improvements, func(i, j int) bool {
		return improvements[i].PointsDelta > improvements[j].PointsDelta
	})

	if len(improvements) > 3 {
		improvements = improvements[:3]
	}

	return improvements
}

func generateInsights(components []ComponentScore, pattern string) []string {
	var insights []string

	switch pattern {
	case "LIQUIDITY_GAP_SAVER":
		insights = append(insights,
			"You're a consistent saver who occasionally breaks FDs for cash flow.",
			"A credit line could eliminate the need to break FDs while building credit.",
			"Consider maintaining a small liquid fund to avoid premature withdrawals.")
	case "LOYAL_SINGLE_BANK":
		insights = append(insights,
			"Your savings discipline is excellent, but keeping everything at one bank may mean missing better rates.",
			"You might be exceeding DICGC insurance limits (₹5L per bank).",
			"Consider spreading across banks for better rate optimization.")
	case "RECOVERING_SAVER":
		insights = append(insights,
			"You've shown resilience by recovering from a financial setback.",
			"Your growth trajectory is positive — keep it up!",
			"Maintaining consistency will further improve your score.")
	case "DISCIPLINED_OPTIMIZER":
		insights = append(insights,
			"Excellent! You demonstrate high financial discipline across all metrics.",
			"You're eligible for premium credit products with favorable terms.",
			"Your FD history makes you a low-risk borrower.")
	default:
		consistency := findComponent("Consistency", components)
		discipline := findComponent("Maturity Discipline", components)
		growth := findComponent("Growth Trajectory", components)
		diversification := findComponent("Diversification", components)

		if consistency.Score < 40 {
			insights = append(insights, "Focus on making regular, consistent deposits to improve your score.")
		}
		if discipline.Score < 40 {
			insights = append(insights, "Avoid premature FD withdrawals — hold to maturity when possible.")
		}
		if growth.Score < 40 {
			insights = append(insights, "Gradually increase your deposit amounts year over year.")
		}
		if diversification.Score < 40 {
			insights = append(insights, "Consider diversifying across different banks and tenure types.")
		}
	}

	if len(insights) == 0 {
		insights = append(insights, "Your FD habits show a balanced approach to savings.")
		insights = append(insights, "Continue maintaining good financial discipline.")
	}

	return insights
}

var _ = sort.Ints