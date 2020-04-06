package gospinner

import (
	"fmt"
	"os"
	"time"
)

var (
	// Default Values
	defaultSpinCharacters = "|/-\\"
	defaultHz             = 15

	// Terminal Control Codes
	terminalBackspace = []byte("\010")
	terminalCursorOn  = []byte("\033[?25h")
	terminalCursorOff = []byte("\033[?25l")
)

type Spinner struct {
	SpinCharacters string
	Hz             int
	CursorVisible  bool
	IsRunning      bool

	// Mostly for testing
	MaxFrames   uint64
	MaxDuration time.Duration

	// Internal counter for frame position
	Frames uint64

	// Internal channel to signal stop
	stopChannel chan bool
}

func (s *Spinner) ensureDefaults() {
	if s.SpinCharacters == "" {
		s.SpinCharacters = defaultSpinCharacters
	}

	if s.Hz < 1 {
		s.Hz = defaultHz
	}

	if s.stopChannel == nil {
		s.stopChannel = make(chan bool, 1)
	}
}

func showCursor(visible bool) {
	if visible {
		os.Stdout.Write(terminalCursorOn)
	} else {
		os.Stdout.Write(terminalCursorOff)
	}
}

func (s *Spinner) Start() {
	s.ensureDefaults()

	sleepTime := time.Millisecond * time.Duration(1000/s.Hz)
	showCursor(s.CursorVisible)
	s.IsRunning = true

	go func() {
		for {
			os.Stdout.Write([]byte{s.SpinCharacters[s.Frames%uint64(len(s.SpinCharacters))]})
			time.Sleep(sleepTime)
			os.Stdout.Write(terminalBackspace)

			if s.MaxFrames > 0 && s.Frames == s.MaxFrames {
				break
			}

			s.Frames++

			select {
			case <-s.stopChannel:
				fmt.Println("got stop")
				s.IsRunning = false
				close(s.stopChannel)
				break
			case <-time.After(sleepTime):
				fmt.Println("no stop")
				continue
			}

		}
	}()
}

func (s *Spinner) Stop() {
	showCursor(true)

	s.stopChannel <- true
}
