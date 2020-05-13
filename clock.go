package main

import (
	"github.com/fogleman/gg"
	"math"
	"time"
)

// Draws an analog clock face filling a square of size pixels
// with top-left corner at x/y
func drawClock(dc *gg.Context, x, y float64, size float64) {
	r := size / 2
	cx, cy := x+r, y+r
	lw := 1.5
	if size < 32 {
		lw = 0.4
	}
	now := time.Now()
	hour, min := now.Hour(), now.Minute()

	// Draw outer circle and fill with color
	dc.SetRGB255(255, 255, 255)
	dc.DrawCircle(cx, cy, r)
	dc.Fill()

	// Draw slightly smaller circle and fill with black
	dc.SetRGB255(0, 0, 0)
	dc.DrawCircle(cx, cy, r*0.95)
	dc.Fill()

	// Draw minute hand
	dc.SetRGB255(0, 0, 127)
	dc.MoveTo(cx, cy)
	dc.SetLineWidth(lw)
	dx, dy := minSecCords(min, cx, cy, r*0.90)
	dc.DrawLine(cx, cy, dx, dy)
	dc.Stroke()

	// Draw hour hand
	dx, dy = hourCords(hour, min, cx, cy, r*0.55)
	dc.SetRGB255(0, 0, 0)
	dc.DrawLine(cx, cy, dx, dy) // first draw in black so we don't mix the hour hand color with minute hand color
	dc.SetRGB255(0, 90, 127)
	dc.DrawLine(cx, cy, dx, dy)
	dc.Stroke()
}

// Given current minute or second time(i.e 30 min, 60 minutes)
// and the radius, returns pair of cords to draw line to
func minSecCords(n int, cx, cy, radius float64) (float64, float64) {
	// converts min/sec to angle and then to radians
	theta := (float64(n)*(360/60) - 90) * (math.Pi / 180)
	x, y := radius*math.Cos(theta), radius*math.Sin(theta)
	return x + cx, y + cy
}

// Given current hour time(i.e. 12, 8) and the radius,
// returns pair of cords to draw line to
func hourCords(hour, min int, cx, cy, radius float64) (float64, float64) {
	// converts hours to angle and then to radians
	theta := ((float64(hour%12)+(float64(min)/60))*(360/12) - 90) * (math.Pi / 180)
	x, y := radius*math.Cos(theta), radius*math.Sin(theta)
	return x + cx, y + cy
}
