package main

import (
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

type splittable interface {
	split(total money.Money) []ParticipantOwed
}

type EvenSplit struct {
	Participants uuid.UUIDs
}
type AdjustmentSplit []struct {
	UserUuid   uuid.UUID
	Adjustment int
}
type PercentSplit []struct {
	UserUuid uuid.UUID
	Percent  int
}

type ParticipantOwed struct {
	UserUuid   uuid.UUID
	AmountOwed int
}

func (s *EvenSplit) split(total money.Money) []ParticipantOwed {

	panic("split not implemented")
}

type Wow struct {
	Total     int
	SplitType splittable
}

func (w *Wow) split(total money.Money) {
	w.SplitType.split(total)
}
