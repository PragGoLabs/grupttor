package grupttor

// Handler is func for specify handler interface
type Handler func(interrupter *Grupttor) error

// InterruptorState represents internal state of interruptor
type InterruptorState string

const (
	// INIT grupttor is initialized but not started
	INIT InterruptorState = "init"
	// WAITING grupttor is started and waiting for signals
	WAITING InterruptorState = "waiting"
	// INTERRUPTING received interruption signal and running interruption handler
	INTERRUPTING InterruptorState = "interrupting"
	// STOPPED received stop signal and running stop handler
	STOPPED InterruptorState = "stopped"
)
