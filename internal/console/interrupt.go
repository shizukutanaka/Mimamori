package console

import (
	"os"
	"sync"
)

var statusControlCExit uint32 = 0xC000013A // STATUS_CONTROL_C_EXIT

func controlCExitCode() int { return int(int32(statusControlCExit)) }

// restoreOnInterrupt runs restore once, either on the first interrupt (then exit) or when stop is called.
func restoreOnInterrupt(interrupts <-chan os.Signal, restore func(), exit func()) (stop func()) {
	var restoreOnce sync.Once
	var stopOnce sync.Once
	done := make(chan struct{})
	finished := make(chan struct{})

	go func() {
		defer close(finished)
		select {
		case <-interrupts:
			restoreOnce.Do(restore)
			exit()
		case <-done:
		}
	}()

	return func() {
		stopOnce.Do(func() { close(done) })
		<-finished
		restoreOnce.Do(restore)
	}
}
