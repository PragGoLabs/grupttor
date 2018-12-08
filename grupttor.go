package grupttor

import (
	"os"
	"time"
)

type Grupttor struct {
	interrupterState InterruptorState
	interruptChannel chan os.Signal
	stopChannel      chan os.Signal

	interruptHandler Handler
	stopHandler      Handler

	interruptHooks []Hook
}

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

func (i *Grupttor) IsInit() bool {
	return i.interrupterState == INIT
}

func (i *Grupttor) AddHook(hook Hook) error {
	if !i.IsInit() {
		return CreateInterruptorWrongStateError("interruptor already waiting for signals, add hooks is not allowed")
	}

	i.interruptHooks = append(i.interruptHooks, hook)
	return nil
}

func (i *Grupttor) GetState() InterruptorState {
	return i.interrupterState
}

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

func (i *Grupttor) IsWaiting() bool {
	return i.interrupterState == WAITING
}

func (i *Grupttor) Interrupt() error {
	if !i.IsWaiting() {
		return CreateInterruptorWrongStateError("Interruptor is not in waiting state, unable to interrupt, skipping")
	}

	i.interruptChannel <- os.Interrupt
	return nil
}

func (i *Grupttor) IsInterrupting() bool {
	return i.interrupterState == INTERRUPTING
}

func (i *Grupttor) Stop() error {
	if !i.IsInterrupting() {
		return CreateInterruptorWrongStateError("Interruptor is not in interrupting state, unable to stop, skipping")
	}

	i.stopChannel <- os.Interrupt
	return nil
}

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
