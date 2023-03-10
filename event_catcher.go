package main

import (
	"github.com/eiannone/keyboard"
	"github.com/tevino/abool"
	"os"
	"os/signal"
	"syscall"
)

// https://opensourcedoc.com/golang-programming/class-object/
type EventCatcher struct {
	windowChangeEvent *abool.AtomicBool
	stopEvent         *abool.AtomicBool
}

func (e *EventCatcher) listenKeystroke() {
	// https://github.com/eiannone/keyboard
	// Exit after pressing any key.
	go func() {
		if err := keyboard.Open(); err != nil {
			panic(err)
		}
		defer func() {
			_ = keyboard.Close()
		}()

		keyboard.GetKey()
		e.stopEvent.Set()
	}()
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
				e.stopEvent.Set()
				break
			case syscall.SIGWINCH:
				e.windowChangeEvent.Set()
			}
		}
	}()
}
