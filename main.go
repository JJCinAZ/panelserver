package main

import (
	"context"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/jjcinaz/panelserver/minicron"
	"github.com/jjcinaz/panelserver/pixelpusher"
	"golang.org/x/image/font/inconsolata"
	"image"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flag"
	"github.com/rs/zerolog"
)

var (
	debugmode    bool
	rows, cols   int
	address      string
	pgmTerminate context.CancelFunc
	logger       zerolog.Logger
)

type calendarData struct {
	upcomingEvent bool
	lines         [2]string
}

type marketData struct {
	lines [2]string
}

type PanelData struct {
	displaymode int
	caldata     calendarData
	mktdata     marketData
}

func main() {
	var (
		ctxExiting context.Context
		paneldata  PanelData
	)
	flag.BoolVar(&debugmode, "debug", false, "Enable verbose output")
	flag.IntVar(&rows, "rows", 0, "Rows on RGB matrix panel")
	flag.IntVar(&cols, "cols", 0, "Cols on RGB matrix panel")
	flag.StringVar(&address, "address", "", "IP address of panel")
	flag.Parse()
	if debugmode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		//		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	ctxExiting, pgmTerminate = context.WithCancel(context.Background())
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt, syscall.SIGTERM:
			pgmTerminate()
		}
	}()
	sched := minicron.NewSchedule()
	// InitWebServer()
	sched.AddJob(time.Minute*5, true, getCalendar, &paneldata.caldata)
	sched.AddJob(time.Minute*5, true, getMarketData, &paneldata.mktdata)
	sched.AddJob(time.Second*5, true, updateRGBmatrix, &paneldata)
	timer := time.NewTicker(time.Second)
MAINLOOP:
	for {
		select {
		case <-timer.C:
			sched.ServiceNextJob()
		case <-ctxExiting.Done():
			break MAINLOOP
		}
	}
	logger.Info().Msg("Exiting normally")
}

func updateRGBmatrix(params ...interface{}) {
	paneldata := params[0].(*PanelData)
	logger.Debug().Int("displaymode", paneldata.displaymode).Msg("Updating RGB Matrix")
	client, err := pixelpusher.NewClient("udp", address, rows, cols)
	if err != nil {
		logger.Error().Err(err)
		return
	}
	defer client.Close()
	switch paneldata.displaymode {
	case 0, 1, 2, 3:
		err = client.SendImage(buildClockImage(paneldata))
		if paneldata.caldata.upcomingEvent {
			paneldata.displaymode++
		} else {
			paneldata.displaymode = 4
		}
	case 4:
		err = client.SendImage(buildMarketImage(paneldata))
		if paneldata.caldata.upcomingEvent {
			paneldata.displaymode = 0
		}
	}
	if err != nil {
		fmt.Println(err)
	}
}

func buildClockImage(paneldata *PanelData) image.Image {
	dc := gg.NewContext(cols, rows)
	dc.SetFontFace(inconsolata.Regular8x16)
	dc.SetRGB255(0, 0, 0)
	dc.Clear()
	drawClock(dc, float64(cols)-32, 0, 32)
	dc.SetRGB255(255, 255, 255)
	s := paneldata.caldata.lines[0]
	if len(s) > 12 {
		s = s[0:12]
	}
	dc.DrawStringAnchored(s, 0, 0, 0, 0.8)
	s = paneldata.caldata.lines[1]
	if len(s) > 12 {
		s = s[0:12]
	}
	dc.DrawStringAnchored(s, 0, 16, 0, 0.8)
	return dc.Image()
}

func buildMarketImage(paneldata *PanelData) image.Image {
	dc := gg.NewContext(cols, rows)
	dc.SetFontFace(inconsolata.Regular8x16)
	dc.SetRGB255(0, 0, 0)
	dc.Clear()
	dc.SetRGB255(255, 255, 255)
	s := paneldata.mktdata.lines[0]
	if len(s) > 16 {
		s = s[0:16]
	}
	dc.DrawStringAnchored(s, 0, 0, 0, 0.8)
	s = paneldata.mktdata.lines[1]
	if len(s) > 16 {
		s = s[0:16]
	}
	dc.DrawStringAnchored(s, 0, 16, 0, 0.8)
	return dc.Image()
}
