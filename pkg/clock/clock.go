package clock

import "time"

type Clock interface {
	Now() time.Time
	NowPointer() *time.Time
}

type RealClock struct {
}

func NewRealClock() Clock {
	return &RealClock{}
}

func (r RealClock) Now() time.Time {
	return time.Now()
}

func (r RealClock) NowPointer() *time.Time {
	now := r.Now()
	return &now
}
