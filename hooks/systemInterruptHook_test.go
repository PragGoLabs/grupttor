package hooks

import (
	"github.com/PragGoLabs/grupttor"
	"github.com/PragGoLabs/grupttor/handlers"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestNewSystemInterruptHook(t *testing.T) {
	testStatusChan := make(chan bool, 1)

	// setup interruptor
	interruptor := grupttor.NewGrupttor(
		handlers.NewWrapHandler(
			func(interrupter *grupttor.Grupttor) error {
				_ = interrupter.Stop()

				return nil
			},
			func(interrupter *grupttor.Grupttor) error {
				// pass false to chan
				testStatusChan <- true
				
				return nil
			},
		),
		[]grupttor.Hook{},
	)

	err := interruptor.AddHook(NewSystemInterruptHook([]os.Signal{syscall.SIGKILL, syscall.SIGTERM}))
	if err != nil {
		t.Fatal(err)
	}

	// for case when interruptor didnt work - test fail
	go func() {
		time.Sleep(5 * time.Second)

		// stop for exit, with error
		testStatusChan <- false
	}()

	// run interruptor and wait
	go interruptor.StartAndWait()

	// system sigterm simulation
	go func() {
		time.Sleep(3 * time.Second)

		// and send interrupt
		err := syscall.Kill(os.Getpid(), syscall.SIGTERM)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// wait until everything run
	testStatus := <-testStatusChan
	if !testStatus {
		t.Fatalf("Interruptor is still running with state %s", interruptor.GetState())
	}
}
