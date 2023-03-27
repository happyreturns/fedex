package clock

import (
	"time"
)

//go:generate mockgen -package=mocks -source=clock.go -destination=mocks/clock.go
type Clock interface {
	Now() time.Time
}

type clock struct{}

func NewClock() *clock {
	return &clock{}
}

func (c *clock) Now() time.Time {
	return time.Now()
}
