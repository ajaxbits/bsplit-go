package main

import (
	"math"
)

type SplitType int

const (
	Even = iota
	Percent
	Adjustment
	Exact
)

var splitTypeName = map[SplitType]string{
	Even:       "even",
	Percent:    "percent",
	Adjustment: "adjustment",
	Exact:      "exact",
}

func Split(amount float64, participants int, splitType SplitType) []int {
	switch splitType {
	case Even:
		amountCents := int(math.Round(amount * 100))
		return evenSplit(amountCents, participants)
	default:
		panic("Not implemented")
	}
}

func evenSplit(amount int, participants int) []int {
	if participants <= 0 {
		return []int{amount}
	}

	baseAmount := amount / participants

	splits := make([]int, participants)
	for i := range splits {
		splits[i] = baseAmount
	}

	totalAssigned := baseAmount * participants
	remainingAmount := amount - totalAssigned

	for i := 0; remainingAmount > 0 && i < participants; i++ {
		splits[i] = splits[i] + 1
		remainingAmount = remainingAmount - 1
	}

	return splits
}
