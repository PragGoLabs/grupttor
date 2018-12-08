package grupttor

// Handle interface specify handle function
type Handle interface {
	// HandleInterrupt handler interrupt signal
	HandleInterrupt(interrupter *Grupttor) error

	// HandleStop handle stop signal
	HandleStop(interrupter *Grupttor) error
}