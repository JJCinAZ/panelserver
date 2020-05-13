package main

import (
	"context"
	"fmt"
	"github.com/jjcinaz/msgraph"
	"github.com/kofoworola/godate"
	"os"
	"sort"
	"time"
)

func getCalendar(params ...interface{}) {
	var (
		err                           error
		c                             *msgraph.Client
		tenantid, clientid, clientkey string
	)
	logger.Debug().Msg("Getting Calendar")
	caldata := params[0].(*calendarData)
	caldata.upcomingEvent = false
	caldata.lines[0] = "See"
	caldata.lines[1] = "receptionist"
	tenantid = os.Getenv("AZURE_TENANTID")
	clientid = os.Getenv("AZURE_CLIENTID")
	clientkey = os.Getenv("AZURE_CLIENTKEY")
	if len(tenantid) == 0 {
		logger.Error().Msg("Missing environment variable AZURE_TENANTID")
		return
	}
	if len(clientid) == 0 {
		logger.Error().Msg("Missing environment variable AZURE_CLIENTID")
		return
	}
	if len(clientkey) == 0 {
		logger.Error().Msg("Missing environment variable AZURE_CLIENTKEY")
		return
	}
	c, err = msgraph.NewKeyClient(context.Background(), tenantid, clientid, clientkey)
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}
	defer c.Close()
	daystart := godate.Create(time.Now())
	events, err := c.GetCalendarView("SES-LgConf@simplybits.com", msgraph.OptionTextMailBody(),
		msgraph.OptionStartDateTime(daystart.Time), msgraph.OptionEndDateTime(daystart.EndOfDay().Time))
	if err != nil {
		logger.Error().Err(err).Send()
		return
	}
	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Native.Before(events[j].Start.Native)
	})
	logger.Debug().Msgf("Got %d events", len(events))
	if len(events) == 0 {
		caldata.upcomingEvent = false
		caldata.lines[0] = "Available"
		caldata.lines[1] = ""
	} else {
		logger.Debug().Msgf("%+v", events[0])
		caldata.upcomingEvent = true
		s := events[0].Start.Native.Local()
		e := events[0].End.Native.Local()
		if s.Add(-5 * time.Minute).After(time.Now()) {
			// Event is upcoming within 5 minutes
			caldata.lines[0] = events[0].Subject
			caldata.lines[1] = fmt.Sprintf("%s…%s", s.Format("3:04"), e.Format("3:04"))
		} else if s.After(time.Now()) {
			caldata.lines[0] = "Available"
			caldata.lines[1] = fmt.Sprintf("until %s", s.Format("3:04"))
		} else {
			caldata.lines[0] = events[0].Subject
			caldata.lines[1] = fmt.Sprintf("until %s", e.Format("3:04"))
		}
	}
	logger.Debug().Strs("lines", caldata.lines[:]).Send()
}
