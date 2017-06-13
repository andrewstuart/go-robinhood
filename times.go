package robinhood

import "time"

// const (
// 	am = 0
// 	pm = 1
// )

// type timeOfDay struct {
// 	Hour, Minute, Meridian int
// }

// func (t timeOfDay) minute() int {
// 	return (t.Meridian*12+t.Hour)*60 + t.Minute
// }

// func tdFromTime(t time.Time) timeOfDay {
// 	return td{time.Hour % 12, time.Minute, time.Hour / 12}
// }

//
const (
	ExtendedOpen    = 4 * 60
	RHExtendedOpen  = 9 * 60
	OpenMinute      = 9*60 + 30
	CloseMinute     = 16 * 60
	RHExtendedClose = 18 * 60
	ExtendedClose   = 20 * 60
)

// MinuteOfDay returns the minute of the day for a given time.Time (hr * 60 +
// min)
func MinuteOfDay(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}

func nyLoc() *time.Location {
	et, _ := time.LoadLocation("America/New_York")
	return et
}

func nyMinute() int {
	return MinuteOfDay(time.Now().In(nyLoc()))
}

func isWeekday(t time.Time) bool {
	wd := t.Weekday()
	return wd != time.Saturday && wd != time.Sunday
}

// IsWeekDay returns whether the given time is a regular
// weekday
func IsWeekDay(t time.Time) bool {
	return isWeekday(time.Now())
}

// NextWeekday returns the next weekday
func NextWeekday() time.Time {
	d := time.Now().AddDate(0, 0, 1)
	for !isWeekday(d) {
		d = d.AddDate(0, 0, 1)
	}
	return d
}

// IsRegularTradingTime returns whether or not the markets are currently open
// for regular trading.
func IsRegularTradingTime() bool {
	now := nyMinute()
	return OpenMinute <= now && now < CloseMinute
}

// IsRobinhoodExtendedTradingTime returns whether or not trades can still be
// placed during the robinhood gold extended trading hours.
func IsRobinhoodExtendedTradingTime() bool {
	now := nyMinute()
	return RHExtendedOpen <= now && now < RHExtendedClose
}

// IsExtendedTradingTime returns whether or not extended hours equity will be
// updated because extended-hours trades may still be allowed in the markets.
func IsExtendedTradingTime() bool {
	now := nyMinute()
	return ExtendedOpen <= now && now < ExtendedClose
}

func nextDayHourDateNY(h, m int) time.Time {
	now := time.Now()

	if h*60+m <= MinuteOfDay(now) {
		now = NextWeekday()
	}

	return time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, nyLoc())
}

// NextMarketOpen returns the time of the next opening bell, when regular
// trading begins.
func NextMarketOpen() time.Time {
	return nextDayHourDateNY(9, 30)
}

// NextMarketExtendedOpen returns the time of the next extended opening time,
// when stock equity may begin to fluctuate again.
func NextMarketExtendedOpen() time.Time {
	return nextDayHourDateNY(4, 00)
}

// NextRobinhoodExtendedOpen returns the time of the next extended
func NextRobinhoodExtendedOpen() time.Time {
	return nextDayHourDateNY(9, 00)
}
