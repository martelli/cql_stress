package ralpe

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jonboulle/clockwork"
)

type ralpe struct {
	rate          int
	num_workers   int
	total_jobs    int
	wg            sync.WaitGroup
	ch            chan bool
	fun           func() error
	clock         clockwork.Clock
	completed     atomic.Uint64
	accum_latency atomic.Uint64
}

func NewRalpe(fun func() error, rate, parallel, total int) *ralpe {
	r := &ralpe{
		fun:         fun,
		rate:        rate,
		num_workers: parallel,
		total_jobs:  total,
		ch:          make(chan bool, rate),
		clock:       clockwork.NewRealClock(),
	}

	for i := range parallel {
		r.wg.Add(1)
		go r.worker(i)
	}
	return r
}

func (r *ralpe) SetClock(c clockwork.Clock) {
	r.clock = c
}

func (r *ralpe) worker(id int) {
	for _ = range r.ch {
		start := time.Now()
		err := r.fun()
		if err != nil {
			slog.Error("Failed query", "error", err.Error())
			return
		}
		latency := time.Now().Sub(start)
		r.completed.Add(1)
		r.accum_latency.Add(uint64(latency))
	}
	r.wg.Done()
}

func (r *ralpe) produce() {

	total := r.total_jobs

	fill := func(c int) int {
		count := 0
	forLoop:
		for _ = range c {
			select {
			case r.ch <- true:
				count++
			default:
				break forLoop
			}
		}
		return count
	}

	for total > 0 {
		amount := r.rate
		if total < r.rate {
			amount = total
		}
		sent := fill(amount)
		<-r.clock.After(time.Second)
		r.Stats()
		total -= sent
	}

	close(r.ch)
}

func (r *ralpe) Start() {
	go r.produce()
}

func (r *ralpe) Wait() {
	r.wg.Wait()
}

func (r *ralpe) Stats() {
	throughput := r.completed.Swap(0)
	accum_latency := r.accum_latency.Swap(0)
	if throughput == 0 {
		return
	}
	avg := time.Duration(accum_latency / throughput)
	slog.Info("Stats:", "throughput", throughput, "avg_latency:", avg)
}
