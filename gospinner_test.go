package gospinner

import (
	"testing"
	"time"
)

func TestMaxFrames(t *testing.T) {
	s := Spinner{}
	s.MaxFrames = 10
	s.Start()

	for {
		if !s.IsRunning {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if s.Frames != s.MaxFrames {
		t.Fatalf("s.Frames(%v) != s.MaxFrames (%v)", s.Frames, s.MaxFrames)
	}

}
