package hooks

import (
	"github.com/PragGoLabs/grupttor"
	"os"
	"os/signal"
)

// SystemInterruptHook contains all signals it will handle
type SystemInterruptHook struct {
	allowedSignals []os.Signal
}

// NewSystemInterruptHook factory for creating new system interrupt with allowed signals
func NewSystemInterruptHook(allowedSignals []os.Signal) SystemInterruptHook {
	return SystemInterruptHook{
		allowedSignals: allowedSignals,
	}
}

// Init will attach on grupttor interrupt signal
func (sih SystemInterruptHook) Init(interrupter *grupttor.Grupttor) {
	// create buffered channel of os signals
	sigChannel := make(chan os.Signal, 1)

	// register system signal notification
	signal.Notify(sigChannel, sih.allowedSignals...)

	select {
	case <-sigChannel:
		// send interrupt
		err := interrupter.Interrupt()

		// there is something wrong in application state
		if err != nil {
			panic(err)
		}
	}
}
