package handlers

import (
	"github.com/PragGoLabs/grupttor"
)

// WrapHandler is default handler, when you just want to pass simple functions
type WrapHandler struct {
	interruptHandler grupttor.Handler
	stopHandler      grupttor.Handler
}

// NewWrapHandler create struct with internal func handlers
func NewWrapHandler(interruptHandler grupttor.Handler, stopHandler grupttor.Handler) WrapHandler {
	return WrapHandler{
		interruptHandler: interruptHandler,
		stopHandler:      stopHandler,
	}
}

// HandleInterrupt handler interrupt signal
func (wh WrapHandler) HandleInterrupt(interruptor *grupttor.Grupttor) error {
	return wh.interruptHandler(interruptor)
}

// HandleStop handle stop signal
func (wh WrapHandler) HandleStop(interruptor *grupttor.Grupttor) error {
	return wh.stopHandler(interruptor)
}
