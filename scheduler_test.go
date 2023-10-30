package scheduler

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSchedulerOne(t *testing.T) {
	var wg sync.WaitGroup

	var counter atomic.Int32

	wg.Add(1)
	s := NewScheduler(
		WithInterval(10*time.Millisecond),
		WithHandler(HandlerFunc(func(t time.Time) {
			defer wg.Done()

			counter.Add(1)
		})),
	)

	go func() {
		err := s.Start()

		assert.NoError(t, err)
	}()

	wg.Wait()

	err := s.Stop()
	assert.NoError(t, err)

	assert.Equal(t, 1, int(counter.Load()))

	time.Sleep(11 * time.Millisecond)
}

func TestSchedulerTwo(t *testing.T) {
	var wg sync.WaitGroup

	var counter atomic.Int32

	wg.Add(2)
	s := NewScheduler(
		WithInterval(10*time.Millisecond),
		WithHandler(HandlerFunc(func(t time.Time) {
			defer wg.Done()

			counter.Add(1)
		})),
	)

	go func() {
		err := s.Start()

		assert.NoError(t, err)
	}()

	wg.Wait()

	err := s.Stop()
	assert.NoError(t, err)

	assert.Equal(t, 2, int(counter.Load()))

	time.Sleep(11 * time.Millisecond)
}
