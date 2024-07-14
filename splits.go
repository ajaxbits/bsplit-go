package main

import (
	"math/rand"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

type splittable interface {
	split(total money.Money) ([]ParticipantOwed, error)
}

type ParticipantOwed struct {
	UserUuid   uuid.UUID
	AmountOwed money.Money
}

func Split[S splittable](total int64, splitType S) ([]ParticipantOwed, error) {
	return splitType.split(*money.New(total, money.USD))
}

// Strategies

type EvenSplit struct {
	Participants uuid.UUIDs
}
type AdjustmentSplit []struct {
	UserUuid   uuid.UUID
	Adjustment int
}
type PercentSplit []struct {
	UserUuid uuid.UUID
	Percent  int64
}

func scrambleSlice[T any](s []T) []T {
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s
}

func (s *EvenSplit) split(total money.Money) ([]ParticipantOwed, error) {
	shares, err := total.Split(len(s.Participants))
	if err != nil {
		return nil, err
	}

	// Make sure it's kind of fair due to round-robin misalignments
	scrambledParticipants := scrambleSlice(s.Participants)

	wow := Zip(scrambledParticipants, shares)

	result := make([]ParticipantOwed, len(s.Participants))
	for i, p := range wow {
		participantUuid, share := p.First, p.Second
		result[i] = ParticipantOwed{
			UserUuid:   participantUuid,
			AmountOwed: *share,
		}
	}

	return result, nil
}

func (s *PercentSplit) split(total money.Money) ([]ParticipantOwed, error) {
	allocations := make([]int, len(*s))
	for i, p := range *s {
		allocations[i] = int(p.Percent)
	}

	shares, err := total.Allocate(allocations...)
	if err != nil {
		return nil, err
	}

	// Make sure it's kind of fair due to round-robin misalignments
	scrambled := scrambleSlice(*s)

	paired := Zip(scrambled, shares)

	result := make([]ParticipantOwed, len(*s))
	for i, p := range paired {
		result[i] = ParticipantOwed{
			UserUuid:   p.First.UserUuid,
			AmountOwed: *p.Second,
		}
	}

	return result, nil
}
