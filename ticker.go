package uneventicker

import (
	"context"
	"sync"
	"time"
)

// Ticker provides uneven ticker events.
type Ticker struct {
	C  <-chan time.Time
	mu sync.Mutex
	cf context.CancelFunc

	durs []time.Duration
}

// New creates `*uneventicker.Ticker` with given durations.
func New(first time.Duration, furthers ...time.Duration) *Ticker {
	durs := make([]time.Duration, len(furthers)+1)
	durs[0] = first
	if len(furthers) > 0 {
		copy(durs[1:], furthers)
	}
	return (&Ticker{
		durs: durs,
	}).start(context.Background())
}

func (ti *Ticker) start(ctx context.Context) *Ticker {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	if ti.cf != nil {
		return ti
	}

	ctx, ti.cf = context.WithCancel(ctx)
	ch := make(chan time.Time)
	ti.C = ch

	next := dursIter(ti.durs)
	tm := time.NewTimer(next())

	go ti.run(ctx, ch, tm, next)

	return ti
}

// dursIter creates a new iterator which repeats the last element when after
// reached to the last.
func dursIter(durs []time.Duration) func() time.Duration {
	return func() time.Duration {
		d := durs[0]
		if len(durs) > 1 {
			durs = durs[1:]
		}
		return d
	}
}

func (ti *Ticker) run(ctx context.Context, ch chan<- time.Time, tm *time.Timer, next func() time.Duration) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case ts := <-tm.C:
			select {
			case ch <- ts:
			default:
			}
			tm.Reset(next())
		}
	}

	if !tm.Stop() {
		<-tm.C
	}
	ti.mu.Lock()
	close(ch)
	ti.mu.Unlock()
}

// Stop stops Ticker to fire any more.
func (ti *Ticker) Stop() bool {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	if ti.cf == nil {
		return false
	}
	ti.cf()
	ti.cf = nil
	return true
}
