package benchmark

import (
	"context"
	"math/rand"
	"time"

	"github.com/vidarandrebo/oncetree/consts"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
)

type Writer struct {
	cfg     *kvsprotos.Configuration
	results []Result
	nodeID  string
	pending chan struct{}
}

func NewWriter(cfg *kvsprotos.Configuration, expectedResults int, nodeID string) *Writer {
	w := &Writer{
		cfg:     cfg,
		nodeID:  nodeID,
		results: make([]Result, 0, expectedResults),
		pending: make(chan struct{}),
	}
	return w
}

func (w *Writer) Run(ctx context.Context) {
	go func() {
	loop:
		for {
			select {
			case <-w.pending:
				startTime := time.Now()
				_, err := w.cfg.Nodes()[0].Write(context.Background(), &kvsprotos.WriteRequest{
					Key:   rand.Int63n(consts.BenchmarkNumKeys),
					Value: rand.Int63n(1000000),
				})
				endTime := time.Now()
				if err != nil {
					// writer has failed, most likely due to a failed replica, we therefore do not set client to ready again
					return
				}
				w.results = append(w.results,
					Result{
						Latency:   endTime.Sub(startTime).Microseconds(),
						Timestamp: endTime.UnixMilli(),
						ID:        w.nodeID,
					})

			case <-ctx.Done():
				break loop
			}
		}
	}()
}

func (w *Writer) Send() bool {
	select {
	case w.pending <- struct{}{}:
		return true
	default:
		return false
	}
}

func (w *Writer) GetResults() []Result {
	return w.results
}
