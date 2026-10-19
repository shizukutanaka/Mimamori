package console

import (
	"os"
	"sync"
	"testing"
	"time"
)

func TestRestoreOnInterrupt(t *testing.T) {
	interrupts := make(chan os.Signal, 1)
	exited := make(chan struct{})
	var mu sync.Mutex
	var calls []string

	stop := restoreOnInterrupt(interrupts, func() {
		mu.Lock()
		calls = append(calls, "restore")
		mu.Unlock()
	}, func() {
		mu.Lock()
		calls = append(calls, "exit")
		mu.Unlock()
		close(exited)
	})
	interrupts <- os.Interrupt

	select {
	case <-exited:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for exit after interrupt")
	}
	stop()

	mu.Lock()
	defer mu.Unlock()
	if len(calls) != 2 || calls[0] != "restore" || calls[1] != "exit" {
		t.Fatalf("calls = %v, want [restore exit]", calls)
	}
}

func TestRestoreOnInterruptStop(t *testing.T) {
	interrupts := make(chan os.Signal, 1)
	exited := make(chan struct{}, 1)
	var mu sync.Mutex
	restoreCalls := 0
	stop := restoreOnInterrupt(interrupts, func() {
		mu.Lock()
		restoreCalls++
		mu.Unlock()
	}, func() {
		exited <- struct{}{}
	})

	stop()
	mu.Lock()
	if restoreCalls != 1 {
		mu.Unlock()
		t.Fatalf("restore called %d times after stop, want once", restoreCalls)
	}
	mu.Unlock()

	interrupts <- os.Interrupt
	select {
	case <-exited:
		t.Fatal("exit called after stop")
	case <-time.After(50 * time.Millisecond):
	}

	mu.Lock()
	defer mu.Unlock()
	if restoreCalls != 1 {
		t.Fatalf("restore called %d times after signal, want once", restoreCalls)
	}
}

func TestRestoreOnInterruptStopTwice(t *testing.T) {
	interrupts := make(chan os.Signal, 1)
	var mu sync.Mutex
	restoreCalls := 0
	stop := restoreOnInterrupt(interrupts, func() {
		mu.Lock()
		restoreCalls++
		mu.Unlock()
	}, func() {})

	stop()
	stop()

	mu.Lock()
	defer mu.Unlock()
	if restoreCalls != 1 {
		t.Fatalf("restore called %d times, want once", restoreCalls)
	}
}

func TestControlCExitCode(t *testing.T) {
	if got := uint32(controlCExitCode()); got != 0xC000013A {
		t.Fatalf("uint32(controlCExitCode()) = %#x, want %#x", got, uint32(0xC000013A))
	}
}
