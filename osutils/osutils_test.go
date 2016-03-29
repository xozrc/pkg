package osutils_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

import (
	"github.com/xozrc/pkg/osutils"
)

var _ = fmt.Print

var signalSlice []syscall.Signal = []syscall.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

var result = 0

func waitSig(t *testing.T, c <-chan os.Signal, sig os.Signal) {
	select {
	case s := <-c:
		if s != sig {
			t.Fatalf("signal was %v, want %v", s, sig)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for %v", sig)
	}
}

func TestInterruptHookRun(t *testing.T) {

	for _, sig := range signalSlice {
		ih := osutils.NewInterruptHooker()

		ih.AddHandler(osutils.InterruptHandlerFunc(firstHandler))
		ih.AddHandler(osutils.InterruptHandlerFunc(secondHandler))
		ih.AddHandler(osutils.InterruptHandlerFunc(thirdHandler))
		ih.RemoveHandler(osutils.InterruptHandlerFunc(thirdHandler))

		c := make(chan os.Signal, 1)
		signal.Notify(c, sig)

		go func() {
			ih.Run()
		}()

		syscall.Kill(syscall.Getpid(), sig)
		waitSig(t, c, sig)
		if result == 3 {
			t.Fatalf("interrupt handlers were called in wrong order")
		}
		if result != 4 {
			t.Fatalf("interrupt handlers were not called properly")
		}
	}
}

func firstHandler() {
	result += 1
}

func secondHandler() {
	result *= 2
}

func thirdHandler() {
	result -= 1
}
