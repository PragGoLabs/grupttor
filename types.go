package grupttor

type Handler func(interrupter *Grupttor)

type InterruptorState string

const (
	INIT         InterruptorState = "init"
	WAITING      InterruptorState = "waiting"
	INTERRUPTING InterruptorState = "interrupting"
	STOPPED      InterruptorState = "stopped"
)
