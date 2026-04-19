package insights

import "fd-credit-score/internal/scoring"

func GeneratePatternNarrative(pattern string, history scoring.DepositHistory) string {
	narratives := map[string]string{
		"LIQUIDITY_GAP_SAVER": "You're a consistent saver who occasionally needs to break FDs for cash flow. A credit line could eliminate this — you'd stop losing FD interest and build credit simultaneously.",
		"LOYAL_SINGLE_BANK": "Your savings discipline is excellent, but keeping everything at one bank means you might be missing better rates elsewhere and exceeding DICGC insurance limits.",
		"RECOVERING_SAVER": "You've shown resilience by recovering from past financial challenges. Your growth pattern is encouraging — keep building on this foundation.",
		"DISCIPLINED_OPTIMIZER": "Outstanding! You represent the ideal FD customer — consistent, disciplined, and strategically smart. You're primed for premium credit products.",
		"CAUTIOUS_DIVERSIFIER": "You're on the right track with your savings journey. Small improvements in tenure management could boost your score significantly.",
		"STANDARD_SAVER": "Your FD habits show a solid foundation. Focus on building consistency and discipline to unlock better credit opportunities.",
	}

	if narrative, ok := narratives[pattern]; ok {
		return narrative
	}
	return "Your financial discipline shows a balanced approach. Keep building good savings habits."
}

func GetProductRecommendations(score int) []scoring.CreditProduct {
	return nil
}

var _ = scoring.DepositHistory{}