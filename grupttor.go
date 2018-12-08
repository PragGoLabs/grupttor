package grupttor

import (
	"os"
	"time"
)

// Grupttor ecosystem struct
type Grupttor struct {
	interrupterState InterruptorState
	interruptChannel chan os.Signal
	stopChannel      chan os.Signal

	interruptHandler Handler
	stopHandler      Handler

	interruptHooks []Hook
}

// NewGrupttor initialize grupttor with all needed handlers passed and with hooks
func NewGrupttor(interruptHandler Handler, stopHandler Handler, hooks []Hook) *Grupttor {
	interrupter := &Grupttor{
		interrupterState: INIT,
		interruptChannel: make(chan os.Signal),
		stopChannel:      make(chan os.Signal),

		interruptHandler: interruptHandler,
		stopHandler:      stopHandler,

		interruptHooks: hooks,
	}

	return interrupter
}

// IsInit check if internal state is INIT
func (i *Grupttor) IsInit() bool {
	return i.interrupterState == INIT
}

// AddHook insert hook to list of existing hooks
// if internal state is not in INIT it will return error
func (i *Grupttor) AddHook(hook Hook) error {
	if !i.IsInit() {
		return CreateInterruptorWrongStateError("interruptor already waiting for signals, add hooks is not allowed")
	}

	i.interruptHooks = append(i.interruptHooks, hook)
	return nil
}

// GetState return internal state
func (i *Grupttor) GetState() InterruptorState {
	return i.interrupterState
}

// StartAndWait change state of interrupt
// and wait for message on interruptChannel
// its main function of grupttor
func (i *Grupttor) StartAndWait() {
	i.changeState(WAITING)

	// init all hooks set
	i.initHooks()

	// wait for interrupt request
	<-i.interruptChannel

	// change state
	i.changeState(INTERRUPTING)

	// do interrupt
	go func() {
		// wait just for the case
		time.Sleep(1 * time.Second)
		// and then interrupt
		i.interruptHandler(i)
	}()

	// wait for interrupt signal
	<-i.stopChannel

	// change state
	i.changeState(STOPPED)

	// run interrupt func
	i.stopHandler(i)
}

// IsWaiting check if internal state is WAITING
func (i *Grupttor) IsWaiting() bool {
	return i.interrupterState == WAITING
}

// Interrupt will pass interrupt to interrupt channel
// if interruptor is not in waiting state it will cause error
func (i *Grupttor) Interrupt() error {
	if !i.IsWaiting() {
		return CreateInterruptorWrongStateError("Interruptor is not in waiting state, unable to interrupt, skipping")
	}

	i.interruptChannel <- os.Interrupt
	return nil
}

// IsInterrupting check if internal state is INTERRUPTING
func (i *Grupttor) IsInterrupting() bool {
	return i.interrupterState == INTERRUPTING
}

// Stop handle sending stop signal to stop channel and cause the end of application
// it expects interruptor is in interrupting state otherwise it cause an error
func (i *Grupttor) Stop() error {
	if !i.IsInterrupting() {
		return CreateInterruptorWrongStateError("Interruptor is not in interrupting state, unable to stop, skipping")
	}

	i.stopChannel <- os.Interrupt
	return nil
}

// IsStopped check if internal state is STOPPED
func (i *Grupttor) IsStopped() bool {
	return i.interrupterState == STOPPED
}

func (i *Grupttor) changeState(state InterruptorState) {
	i.interrupterState = state
}

func (i *Grupttor) initHooks() {
	for _, hook := range i.interruptHooks {
		// wrap every hooks with extra goroutine
		go hook.Init(i)
	}
}
