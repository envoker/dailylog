package dailylog

import (
	"sync"
	"time"
)

type rotator interface {
	Rotate() error
}

func loopRotate(wg *sync.WaitGroup, quit chan struct{}, interval int, r rotator) {

	defer wg.Done()

	r.Rotate()

	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			r.Rotate()
		}
	}
}

type flusher interface {
	Flush() error
}

func loopFlush(wg *sync.WaitGroup, quit chan struct{}, interval int, f flusher) {

	defer wg.Done()

	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			f.Flush()
		}
	}
}
