package hooks

import (
	"github.com/PragGoLabs/grupttor"
	"time"
)

type TimedInterruptHook struct {
	tick time.Duration
}

func NewTimedInterruptHook(tick time.Duration) TimedInterruptHook {
	return TimedInterruptHook{
		tick: tick,
	}
}

func (tih TimedInterruptHook) Init(interrupter *grupttor.Grupttor) {
	ticker := time.NewTimer(tih.tick)

	// waiting for shot
	<-ticker.C

	err := interrupter.Interrupt()
	// there is something wrong in application state
	if err != nil {
		panic(err)
	}
}
