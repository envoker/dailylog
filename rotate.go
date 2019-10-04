package dailylog

import (
	"sync"
	"time"
)

type rotator interface {
	Rotate() error
}

func rotateWorker(wg *sync.WaitGroup, quit <-chan struct{}, d time.Duration, r rotator) {

	defer wg.Done()

	r.Rotate()

	ticker := time.NewTicker(d)
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
