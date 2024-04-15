package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type MetricData struct {
	Method    string
	Duration  time.Duration
	Timestamp time.Time
}

type Metrics struct {
	data   []MetricData
	lock   sync.Mutex
	ticker *time.Ticker
}

func NewMetrics(ctx context.Context, flushInterval time.Duration) *Metrics {
	m := &Metrics{
		ticker: time.NewTicker(flushInterval),
	}
	go m.run(ctx)
	return m
}

func (m *Metrics) Record(method string, duration time.Duration, err error) {
	m.lock.Lock()
	m.data = append(m.data, MetricData{
		Method:    method,
		Duration:  duration,
		Timestamp: time.Now(),
	})
	m.lock.Unlock()
}

func (m *Metrics) run(ctx context.Context) {
	select {
	case <-m.ticker.C:
		m.send()  // Send metrics
		m.flush() // Then flush
	case <-ctx.Done():
		m.ticker.Stop()
		fmt.Println("the ticker has been stopped")

		return
	}
}

func (m *Metrics) flush() {
	m.lock.Lock()
	m.data = []MetricData{} // Reset the slice after copying
	m.lock.Unlock()
}

func (m *Metrics) send() {
	fmt.Printf("sending %d metrics\n", len(m.data))

	// Send this m.data to a telemetry system
}
