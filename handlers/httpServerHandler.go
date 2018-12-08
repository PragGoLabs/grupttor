package handlers

import (
	"context"
	"fmt"
	"github.com/PragGoLabs/grupttor"
	"net/http"
	"sync"
	"time"
)

// DefaultDuration default wait duration until it ends
var DefaultDuration = 8 * time.Second

// HTTPServerHandler it handle interruption to http server
type HTTPServerHandler struct {
	httpServer         *http.Server
	waitGroup          *sync.WaitGroup
	durationOfShutdown *time.Duration
}

// NewHTPPServerHandler create http server handler, for shutdown the http server
func NewHTPPServerHandler(server *http.Server, waitGroup *sync.WaitGroup, duration *time.Duration) HTTPServerHandler {
	if duration == nil {
		duration = &DefaultDuration
	}

	return HTTPServerHandler{
		httpServer:         server,
		waitGroup:          waitGroup,
		durationOfShutdown: duration,
	}
}

// HandleInterrupt handler interrupt signal
func (hsh HTTPServerHandler) HandleInterrupt(interruptor *grupttor.Grupttor) error {
	ctx, cancel := context.WithTimeout(context.Background(), *hsh.durationOfShutdown)
	defer cancel()

	// shutdown the http server
	err := hsh.httpServer.Shutdown(ctx)
	if err != nil {
		return CreateUnableToShutdownHTTPServer(err.Error())
	}
	fmt.Println(hsh.waitGroup)
	// wait until context done
	<-ctx.Done()
	// mark http gr as done
	hsh.waitGroup.Done()

	// and run stop
	_ = interruptor.Stop()

	return nil
}

// HandleStop handle stop signal
func (hsh HTTPServerHandler) HandleStop(interruptor *grupttor.Grupttor) error {
	// there is no need to stop, shutdown close the chan
	hsh.waitGroup.Done()

	return nil
}
