package grupttor

import (
	"testing"
	"time"
)

type PlaceboHook struct{}

func (ph PlaceboHook) Init(grupttor *Grupttor) {}

func TestGrupttor_GetState_INIT(t *testing.T) {
	interuptter := NewGrupttor(
		func(interrupter *Grupttor) {

		},
		func(interrupter *Grupttor) {

		},
		[]Hook{},
	)

	if interuptter.GetState() != INIT {
		t.Fatal("Interrupter is not in init state")
	}
}

func TestGrupttor_GetState_WAITING(t *testing.T) {
	interuptter := createMockInterrupter()

	go interuptter.StartAndWait()

	time.Sleep(2 * time.Second)

	if interuptter.GetState() != WAITING {
		t.Fatal("Interrupter is not in waiting state")
	}
}

func TestGrupttor_AddHook(t *testing.T) {
	interrupter := createMockInterrupter()

	err := interrupter.AddHook(PlaceboHook{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGrupttor_AddHook_Failed(t *testing.T) {
	interrupter := createMockInterrupter()

	go interrupter.StartAndWait()

	time.Sleep(2 * time.Second)

	err := interrupter.AddHook(PlaceboHook{})
	if err == nil {
		t.Fatal("There is error, unable to add hook after interruptor is in waiting state")
	}
}

func createMockInterrupter() *Grupttor {
	return NewGrupttor(
		func(interrupter *Grupttor) {
		},
		func(interrupter *Grupttor) {
		},
		[]Hook{},
	)
}
