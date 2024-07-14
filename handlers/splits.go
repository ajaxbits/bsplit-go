package handlers

import (
	"net/http"
	"strconv"

	"ajaxbits.com/bsplit/splits"
	"ajaxbits.com/bsplit/views"
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
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

	participantUuids := make(uuid.UUIDs, 0, participants)
	for range participants {
		participantUuids = append(participantUuids, uuid.New())
	}

	splits, err := splits.Split(money.NewFromFloat(total, money.USD), &splits.EvenSplit{
		Participants: participantUuids,
	})
	if err != nil {
		c.Logger().Errorf("split failed: %+v", err)
		return c.String(http.StatusInternalServerError, "unable to create split, not enough participants")
	}

	c.Logger().Infof("split: %+v", splits)

	result := make(map[string]string)
	for u, s := range splits {
		result[u.String()] = s.Display()
	}
	return views.Render(c, http.StatusOK, views.Result(result))
}
