package main

import (
	"math/rand"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

type splittable interface {
	split(total money.Money) (ParticipantShares, error)
}

type ParticipantShares map[uuid.UUID]*money.Money

func Split[S splittable](total int64, splitType S) (ParticipantShares, error) {
	return splitType.split(*money.New(total, money.USD))
}

// Strategies

type EvenSplit struct {
	Participants uuid.UUIDs
}
type AdjustmentSplit []struct {
	UserUuid   uuid.UUID
	Adjustment int64
}
type PercentSplit []struct {
	UserUuid uuid.UUID
	Percent  int64
}

func scrambleSlice[T any](s []T) []T {
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s
}

func (s *EvenSplit) split(total money.Money) (ParticipantShares, error) {
	shares, err := total.Split(len(s.Participants))
	if err != nil {
		return nil, err
	}

	// Make sure it's kind of fair due to round-robin misalignments
	scrambledParticipants := scrambleSlice(s.Participants)

	wow := Zip(scrambledParticipants, shares)

	result := make(ParticipantShares)
	for _, p := range wow {
		participantUuid, share := p.First, p.Second
		result[participantUuid] = share
	}

	return result, nil
}

func (s *PercentSplit) split(total money.Money) (ParticipantShares, error) {
	allocations := make([]int, len(*s))
	for i, p := range *s {
		allocations[i] = int(p.Percent)
	}

	shares, err := total.Allocate(allocations...)
	if err != nil {
		return nil, err
	}

	// Make sure it's kind of fair due to round-robin misalignments
	participantsScrambled := scrambleSlice(*s)

	result := make(ParticipantShares, len(*s))
	for _, p := range Zip(participantsScrambled, shares) {
		participantUuid, share := p.First.UserUuid, p.Second
		result[participantUuid] = share
	}

	return result, nil
}

func (s *AdjustmentSplit) split(total money.Money) (ParticipantShares, error) {
	var totalAdjustment int64
	for _, p := range *s {
		totalAdjustment += p.Adjustment
	}

	commonShare := (total.Amount() - totalAdjustment) / int64(len(*s))

	scrambledParticipants := scrambleSlice(*s)
	adjustmentRatios := make([]int, len(*s))
	for i, p := range scrambledParticipants {
		adjustedShare := p.Adjustment + commonShare
		adjustmentRatios[i] = int(adjustedShare)
	}

	shares, err := total.Allocate(adjustmentRatios...)
	if err != nil {
		return nil, err
	}

	result := make(ParticipantShares, len(*s))
	for _, p := range Zip(scrambledParticipants, shares) {
		participantUuid, share := p.First.UserUuid, p.Second
		result[participantUuid] = share
	}

	return result, nil
}
