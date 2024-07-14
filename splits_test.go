package main

import (
	"sort"
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

func TestEvenSplit(t *testing.T) {
	uuid1, uuid2, uuid3 := uuid.New(), uuid.New(), uuid.New()
	splitType := EvenSplit{
		Participants: uuid.UUIDs{uuid1, uuid2, uuid3},
	}
	total := money.New(100, money.USD)

	gotFull, err := splitType.split(*total)
	if err != nil {
		t.FailNow()
	}

	got := make([]int, 0, len(gotFull))
	for _, s := range gotFull {
		got = append(got, int(s.Amount()))
	}
	want := []int{34, 33, 33}

	sort.Ints(got)
	sort.Ints(want)
	for _, p := range Zip(got, want) {
		if p.First != p.Second {
			t.Errorf("got %v, wanted %v", p.First, p.Second)
		}
	}
}
