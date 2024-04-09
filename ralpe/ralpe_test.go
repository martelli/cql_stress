package ralpe

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
)

func TestBucket(t *testing.T) {

	const (
		PARALLEL = 10
		RATE     = 1000
		LOOP     = 20
	)

	var wg sync.WaitGroup

	counter := atomic.Uint32{}
	fkClock := clockwork.NewFakeClock()

	f := func() error {
		counter.Add(1)
		wg.Done()
		return nil
	}
	wg.Add(RATE)
	r := NewRalpe(f, RATE, PARALLEL, LOOP*RATE)
	r.SetClock(fkClock)
	r.Start()

	for _ = range LOOP {
		wg.Wait()
		wg.Add(RATE)
		fkClock.BlockUntil(1)
		fkClock.Advance(1000 * time.Millisecond)
	}

	r.Wait()

	if counter.Load() != LOOP*RATE {
		t.Errorf(">> %d != %d", LOOP*RATE, counter.Load())
	}

}
