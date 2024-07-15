package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ajaxbits.com/bsplit/splits"
	"ajaxbits.com/bsplit/views"
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func SplitHandler(c echo.Context) error {
	formData, err := c.FormParams()
	if err != nil {
		c.Logger().Errorf("invalid total amount: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create split")
	}

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

	participantUuids := make(uuid.UUIDs, 0, participants)
	for range participants {
		participantUuids = append(participantUuids, uuid.New())
	}

	var shares splits.ParticipantShares
	var splitErr error
	switch formData["splitType"][0] {
	case "evenSplit":
		shares, splitErr = splits.Split(money.NewFromFloat(total, money.USD), &splits.EvenSplit{
			Participants: participantUuids,
		})
	case "percentSplit":
		data := make(splits.PercentSplit, 0, len(participantUuids))
		pairs := splits.Zip(participantUuids, strings.Split(formData["percents"][0], ","))
		for _, p := range pairs {
			user, percentStr := p.First, p.Second
			percent, err := strconv.Atoi(percentStr)
			splitErr = err

			percentSplit := struct {
				UserUuid uuid.UUID
				Percent  int64
			}{
				UserUuid: user,
				Percent:  int64(percent)}
			data = append(data, percentSplit)
		}
		shares, splitErr = splits.Split(money.NewFromFloat(total, money.USD), &data)
	}
	if splitErr != nil {
		c.Logger().Errorf("split failed: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create split, not enough participants")
	}

	for u, s := range shares {
		fmt.Println(fmt.Sprintf("user: %s, amount: %s", u, s.Display()))
	}

	result := make(map[string]string)
	for u, s := range shares {
		result[u.String()] = s.Display()
	}
	return views.Render(c, http.StatusOK, views.Result(result))
}
