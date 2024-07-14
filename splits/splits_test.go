package splits

import (
	"reflect"
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

	gotFull, err := splitType.split(total)
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

func TestAdjustmentSplit(t *testing.T) {
	uuid1, uuid2, uuid3 := uuid.New(), uuid.New(), uuid.New()
	splitType := AdjustmentSplit{
		{UserUuid: uuid1, Adjustment: 49},
		{UserUuid: uuid2, Adjustment: 18},
		{UserUuid: uuid3, Adjustment: 0},
	}
	total := money.New(100, money.USD)

	gotRaw, err := splitType.split(total)
	if err != nil {
		t.FailNow()
	}

	got := make(map[uuid.UUID]int64)
	for u, s := range gotRaw {
		got[u] = s.Amount()
	}
	want := map[uuid.UUID]int64{uuid1: 60, uuid2: 29, uuid3: 11}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, wanted %+v", got, want)
	}
}
