package hooks

import (
	"github.com/PragGoLabs/grupttor"
	"time"
)

// TimedInterruptHook contains duration which specify when the interrupt
// shot will happened
type TimedInterruptHook struct {
	tick time.Duration
}

// NewTimedInterruptHook will create timed interrupt task
func NewTimedInterruptHook(tick time.Duration) TimedInterruptHook {
	return TimedInterruptHook{
		tick: tick,
	}
}

// Init will attach on grupttor interrupt signal
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
