package robinhood

import "time"

const (
	ExtendedOpen    = 4 * 60
	RHExtendedOpen  = 9 * 60
	OpenMinute      = 9*60 + 30
	CloseMinute     = 16 * 60
	RHExtendedClose = 18 * 60
	ExtendedClose   = 20 * 60
)

func MinuteOfDay(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}

func nyMinute() int {
	et, _ := time.LoadLocation("America/New_York")
	return MinuteOfDay(time.Now().In(et))
}

// IsRegularTradingTime returns whether or not the markets are currently open
// for regular trading
func IsRegularTradingTime() bool {
	now := nyMinute()
	return OpenMinute <= now && now <= CloseMinute
}

// IsRobinhoodExtendedTradingTime returns whether or not trades can still be
// placed during the robinhood gold extended trading hours
func IsRobinhoodExtendedTradingTime() bool {
	now := nyMinute()
	return RHExtendedOpen <= now && now <= RHExtendedClose
}

// IsExtendedTradingTime returns whether or not extended hours equity will be
// updated because extended-hours trades may still be allowed in the markets.
func IsExtendedTradingTime() bool {
	now := nyMinute()
	return ExtendedOpen <= now && now <= ExtendedClose
}
