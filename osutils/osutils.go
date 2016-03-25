package osutils

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type InterruptHandler interface {
	Run()
}

type InterruptHandlerFunc func()

func (ihf InterruptHandlerFunc) Run() {
	ihf()
}

type InterruptHooker interface {
	AddHandler(ih InterruptHandler)
	Run()
	RemoveHandler(ih InterruptHandler)
}

type interruptHook struct {
	mu       sync.Mutex
	handlers []InterruptHandler
}

func (ih *interruptHook) AddHandler(h InterruptHandler) {
	ih.mu.Lock()
	ih.handlers = append(ih.handlers, h)
	ih.mu.Unlock()
}

func (ih *interruptHook) RemoveHandler(h InterruptHandler) {
	ih.mu.Lock()
	index := 0
	for i, th := range ih.handlers {
		if th == h {
			index = i
		}
	}
	end := 0
	if index > 0 {
		end = index - 1
	}

	ih.handlers = append(ih.handlers[:end], ih.handlers[index:]...)

	ih.mu.Unlock()
}

func (ih *interruptHook) Run() {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//bolck until interrupt signal
	sig := <-notifier

	log.Printf("receive %v signal", sig)

	ih.mu.Lock()
	for _, h := range ih.handlers {
		h.Run()
	}

	ih.mu.Unlock()
	signal.Stop(notifier)

}

func NewInterruptHooker() InterruptHooker {
	return &interruptHook{}
}
