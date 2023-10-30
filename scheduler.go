package scheduler

import (
	"sync/atomic"
	"time"
)

type Handler interface {
	Process(time.Time)
}

type HandlerFunc func(t time.Time)

func (f HandlerFunc) Process(t time.Time) {
	f(t)
}

type Option func(*Scheduler)

func WithInterval(interval time.Duration) Option {
	return func(s *Scheduler) {
		s.interval = interval
	}
}

func WithHandler(handler Handler) Option {
	return func(s *Scheduler) {
		s.handler = handler
	}
}

type Scheduler struct {
	interval time.Duration
	handler  Handler
	shutdown atomic.Bool
}

func NewScheduler(opts ...Option) *Scheduler {
	s := &Scheduler{
		interval: 2 * time.Minute,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Scheduler) process(t time.Time) {
	s.handler.Process(t)
}

func (s *Scheduler) Start() error {
	for {
		if s.shutdown.Load() {
			return nil
		}

		last := time.Now()

		s.process(last)

		elapsed := time.Since(last)
		if elapsed < s.interval {
			time.Sleep(s.interval - elapsed)
		}
	}
}

func (s *Scheduler) Stop() error {
	s.shutdown.Store(true)

	return nil
}
