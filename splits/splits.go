package splits

import (
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

type splittable interface {
	split(total *money.Money) (ParticipantShares, error)
}

type ParticipantShares map[uuid.UUID]*money.Money

func Split[S splittable](total *money.Money, splitType S) (ParticipantShares, error) {
	return splitType.split(total)
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

func (s *EvenSplit) split(total *money.Money) (ParticipantShares, error) {
	shares, err := total.Split(len(s.Participants))
	if err != nil {
		return nil, err
	}

	result := make(ParticipantShares)
	for _, p := range Zip(s.Participants, shares) {
		participantUuid, share := p.First, p.Second
		result[participantUuid] = share
	}

	return result, nil
}

func (s *PercentSplit) split(total *money.Money) (ParticipantShares, error) {
	allocations := make([]int, len(*s))
	for i, p := range *s {
		allocations[i] = int(p.Percent)
	}

	shares, err := total.Allocate(allocations...)
	if err != nil {
		return nil, err
	}

	result := make(ParticipantShares, len(*s))
	for _, p := range Zip(*s, shares) {
		participantUuid, share := p.First.UserUuid, p.Second
		result[participantUuid] = share
	}

	return result, nil
}

func (s *AdjustmentSplit) split(total *money.Money) (ParticipantShares, error) {
	var totalAdjustment int64
	for _, p := range *s {
		totalAdjustment += p.Adjustment
	}

	commonShare := (total.Amount() - totalAdjustment) / int64(len(*s))

	adjustmentRatios := make([]int, len(*s))
	for i, p := range *s {
		adjustedShare := p.Adjustment + commonShare
		adjustmentRatios[i] = int(adjustedShare)
	}

	shares, err := total.Allocate(adjustmentRatios...)
	if err != nil {
		return nil, err
	}

	result := make(ParticipantShares, len(*s))
	for _, p := range Zip(*s, shares) {
		participantUuid, share := p.First.UserUuid, p.Second
		result[participantUuid] = share
	}

	return result, nil
}
