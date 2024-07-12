package handlers

import (
	"math"
	"net/http"
	"strconv"

	"ajaxbits.com/bsplit/views"
	"github.com/labstack/echo/v4"
)

func SplitHandler(c echo.Context) error {
	totalStr, participantsStr := c.FormValue("total"), c.FormValue("participants")
	total, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		c.Logger().Errorf("invalid total amount: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create split")
	}

	participants, err := strconv.Atoi(participantsStr)
	if err != nil {
		c.Logger().Errorf("error parsing participants: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create split")
	} else if participants < 2 {
		c.Logger().Errorf("not enough participants: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create split, not enough participants")
	}

	split := Split(total, participants, Even)

	c.Logger().Infof("split: %+v", split)
	return views.Render(c, http.StatusOK, views.Result(participants, split))
}

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
