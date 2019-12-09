package uneventicker

import (
	"context"
	"testing"
	"time"
)

func timeoutContext(n int, durs []time.Duration) (context.Context, context.CancelFunc) {
	iter := dursIter(durs)
	var sum time.Duration
	for i := 0; i < n; i++ {
		sum += iter()
	}
	return context.WithTimeout(context.Background(), sum*2)
}

func testTicker(t *testing.T, n int, durs ...time.Duration) {
	t.Helper()

	iter := dursIter(durs)
	iter2 := dursIter(durs)
	_ = iter2()

	ctx, cancel := timeoutContext(n, durs)
	defer cancel()
	base := time.Now()
	ti := New(durs[0], durs[1:]...)
	defer ti.Stop()
	for i := 0; i < n; i++ {
		next := iter()
		limit := next + iter2()/2
		select {
		case <-ctx.Done():
			t.Errorf("timeout: i=%d", i)
			return
		case now := <-ti.C:
			d := now.Sub(base)
			base = now
			if d < next {
				t.Errorf("ticker fires too early: i=%d expect=%s actual=%s", i, next, d)
				return
			}
			if d >= limit {
				t.Errorf("ticker fires too late: i=%d expect=%s actual=%s", i, limit, d)
				return
			}
		}
	}
}

func TestTicker(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		testTicker(t, 5, 50*time.Millisecond)
	})
	t.Run("two", func(t *testing.T) {
		testTicker(t, 5, 45*time.Millisecond, 50*time.Millisecond)
	})
	t.Run("three", func(t *testing.T) {
		testTicker(t, 8, 45*time.Millisecond, 100*time.Millisecond,
			50*time.Millisecond)
	})
}
