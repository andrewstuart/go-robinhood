package roho

import "time"

// Common constants for hours and minutes from midnight at which market events
// occur.
const (
	HrExtendedOpen    = 4
	HrRHExtendedOpen  = 9
	HrClose           = 12 + 4
	HrRHExtendedClose = 12 + 6
	HrExtendedClose   = 12 + 8

	MinExtendedOpen    = HrExtendedOpen * 60
	MinRHExtendedOpen  = HrRHExtendedOpen * 60
	MinOpen            = 9*60 + 30
	MinClose           = HrClose * 60
	MinRHExtendedClose = HrRHExtendedClose * 60
	MinExtendedClose   = HrExtendedClose * 60
)

// MinuteOfDay returns the minute of the day for a given time.Time (hr * 60 +
// min).
func MinuteOfDay(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}

// nyLoc returns the *time.Location of New_York.
func nyLoc() *time.Location {
	et, _ := time.LoadLocation("America/New_York")
	return et
}

// nyMinute returns the current minute after midnight in New_York.
func nyMinute() int {
	return MinuteOfDay(time.Now().In(nyLoc()))
}

// isWeekday returns whether or not the given time.Time is a weekday.
func isWeekday(t time.Time) bool {
	wd := t.Weekday()
	return wd != time.Saturday && wd != time.Sunday
}

// IsWeekDay returns whether the given time is a regular
// weekday
func IsWeekDay(t time.Time) bool {
	return isWeekday(time.Now())
}

// NextWeekday returns the next weekday.
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
	return MinOpen <= now && now < MinClose
}

// IsRobinhoodExtendedTradingTime returns whether or not trades can still be
// placed during the robinhood gold extended trading hours.
func IsRobinhoodExtendedTradingTime() bool {
	now := nyMinute()
	return MinRHExtendedOpen <= now && now < MinRHExtendedClose
}

// IsExtendedTradingTime returns whether or not extended hours equity will be
// updated because extended-hours trades may still be allowed in the markets.
func IsExtendedTradingTime() bool {
	now := nyMinute()
	return MinExtendedOpen <= now && now < MinExtendedClose
}

// nextWeekdayHourMinuteNY returns the time.Time of the next h/m occurrence on a weekday in New York.
func nextWeekdayHourMinuteNY(h, m int) time.Time {
	now := time.Now()

	if h*60+m <= MinuteOfDay(now) {
		now = NextWeekday()
	}

	return time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, nyLoc())
}

// NextMarketOpen returns the time of the next opening bell, when regular
// trading begins.
func NextMarketOpen() time.Time {
	return nextWeekdayHourMinuteNY(9, 30)
}

// NextMarketExtendedOpen returns the time of the next extended opening time,
// when stock equity may begin to fluctuate again.
func NextMarketExtendedOpen() time.Time {
	return nextWeekdayHourMinuteNY(HrExtendedOpen, 00)
}

// NextRobinhoodExtendedOpen returns the time of the next robinhood extended
// opening time, when robinhood users can make trades.
func NextRobinhoodExtendedOpen() time.Time {
	return nextWeekdayHourMinuteNY(HrRHExtendedOpen, 00)
}

// NextMarketClose returns the time of the next market close.
func NextMarketClose() time.Time {
	return nextWeekdayHourMinuteNY(HrClose, 00)
}

// NextRobinhoodExtendedClose returns the time of the next robinhood extended
// closing time, when robinhood users must place their last extended-hours
// trade.
func NextRobinhoodExtendedClose() time.Time {
	return nextWeekdayHourMinuteNY(HrRHExtendedClose, 00)
}

// NextMarketExtendedClose returns the time of the next extended market close,
// when stock equity numbers will stop being updated until the next extended
// open.
func NextMarketExtendedClose() time.Time {
	return nextWeekdayHourMinuteNY(HrRHExtendedClose, 00)
}
