package main

import (
	"bufio"
	"github.com/tevino/abool"
	"os"
	"os/signal"
	"syscall"
)

// https://opensourcedoc.com/golang-programming/class-object/
type EventCatcher struct {
	windowChange *abool.AtomicBool
	stop         *abool.AtomicBool
}

func (e *EventCatcher) listenEnter() {
	// https://stackoverflow.com/questions/54422309/how-to-catch-keypress-without-enter-in-golang-loop
	ch := make(chan string)
	go func(ch chan string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, _ := reader.ReadString('\n')
			ch <- s
		}
	}(ch)
	<-ch
	e.stop.Set()
}

func (e *EventCatcher) listenSignal() {
	// https://stackoverflow.com/questions/18106749/golang-catch-signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGWINCH,
	)
	go func() {
		for sig := range sigc {
			switch sig {
			case syscall.SIGINT:
				e.stop.Set()
                break;
			case syscall.SIGWINCH:
				e.windowChange.Set()
			}
		}
	}()
}
