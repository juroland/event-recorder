package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type EventRecorderState int

const (
	stopped EventRecorderState = iota
	started
)

type EventRecorder struct {
	start     time.Time
	end       time.Time
	state     EventRecorderState
	logWriter *bufio.Writer
}

func NewEventRecorder(w *bufio.Writer) *EventRecorder {
	e := new(EventRecorder)
	e.state = stopped
	e.start = time.Now().UTC()
	e.end = time.Now().UTC()
	e.logWriter = w
	return e
}

func (e *EventRecorder) leap() {
	if e.state == stopped {
		e.start = time.Now().UTC()
		e.state = started
		fmt.Printf("%v -> ", e.start)
	} else {
		e.end = time.Now().UTC()
		e.state = stopped
		fmt.Printf("%v\n", e.end)
		e.record()
	}
}

func (e *EventRecorder) record() {
	fmt.Fprintf(e.logWriter, "%v, %v\n", e.start.Format(time.RFC3339), e.end.Format(time.RFC3339))
}

func main() {
	if len(os.Args) != 2 {
		usage(os.Stderr)
		os.Exit(2)
	}

	if os.Args[1] == "--help" {
		usage(os.Stdout)
		return
	}

	logFileName := os.Args[1]
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	eventRecorder := NewEventRecorder(w)

	err = termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeySpace {
				eventRecorder.leap()
			} else if ev.Key == termbox.KeyCtrlD {
				return
			}
		}
	}
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage: event-recorder filename")
}
