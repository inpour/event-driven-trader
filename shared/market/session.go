package market

import "time"

type Session struct {
	StartH, StartM, StopH, StopM int
	Reverse                      bool // if true: session time will be outside the given range
}

func (s *Session) In(unixTime int64) bool {
	t := time.Unix(unixTime, 0).Add(-3*time.Hour - 30*time.Minute)

	if t.Hour() >= s.StartH && t.Hour() < s.StopH || t.Hour() == s.StopH && t.Minute() <= s.StopM {
		return !s.Reverse // unixTime in session range
	}
	return s.Reverse // unixTime not in session range
}
