package main

import (
	// "log"
	"math"
)

func toNearestCent(val float64) float64 {
	return math.Round(val*100) / 100
}

func split(amount float64, participants int, splitType SplitType) []float64 {
	switch splitType {
	case Even:
		return evenSplit(amount, participants)
	default:
		panic("Not implemented")
	}
}

func evenSplit(amount float64, participants int) []float64 {
	if participants <= 0 {
		return []float64{amount}
	}

	baseAmount := toNearestCent(amount / float64(participants))
	splits := make([]float64, participants)

	for i := range splits {
		splits[i] = baseAmount
	}

	totalAssigned := baseAmount * float64(participants)

	remainingAmount := toNearestCent(amount - totalAssigned)

	for i := 0; remainingAmount > 0 && i < participants; i++ {
		splits[i] = toNearestCent(splits[i] + 0.01)
		remainingAmount = toNearestCent(remainingAmount - 0.01)
	}

	return splits
}
