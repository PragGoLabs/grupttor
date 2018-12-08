package hooks

import (
	"github.com/PragGoLabs/grupttor"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestNewTimedInterruptHook(t *testing.T) {
	testStatusChan := make(chan bool, 1)

	// setup interruptor
	interruptor := grupttor.NewGrupttor(
		func(interrupter *grupttor.Grupttor) {
			_ = interrupter.Stop()
		},
		func(interrupter *grupttor.Grupttor) {
			// pass false to chan
			testStatusChan <- true
		},
		[]grupttor.Hook{},
	)

	err := interruptor.AddHook(NewTimedInterruptHook(
		5 * time.Second,
	))

	if err != nil {
		t.Fatal(err)
	}

	// for case when interruptor didnt work - test fail
	go func() {
		time.Sleep(10 * time.Second)

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
