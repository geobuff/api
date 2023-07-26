package utils

func CalculateXPIncrease(score, maxScore int) int {
	percent := (float64(score) / float64(maxScore)) * 100
	if percent > 75 {
		multiplier := 1 + percent/100
		return int(float64((score / 2)) * multiplier)
	}
	return score / 2
}
