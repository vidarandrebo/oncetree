package benchmark

import (
	"context"
	"math/rand"
	"time"

	"github.com/vidarandrebo/oncetree/consts"
	kvsprotos "github.com/vidarandrebo/oncetree/protos/keyvaluestorage"
)

type Reader struct {
	cfg     *kvsprotos.Configuration
	results []Result
	nodeID  string
	pending chan struct{}
}

func NewReader(cfg *kvsprotos.Configuration, expectedResults int, nodeID string) *Reader {
	r := &Reader{
		cfg:     cfg,
		nodeID:  nodeID,
		results: make([]Result, 0, expectedResults),
		pending: make(chan struct{}),
	}
	return r
}

func (r *Reader) Run(ctx context.Context) {
	go func() {
	loop:
		for {
			select {
			case <-r.pending:
				startTime := time.Now()
				_, err := r.cfg.Nodes()[0].Read(context.Background(), &kvsprotos.ReadRequest{
					Key: rand.Int63n(consts.BenchmarkNumKeys),
				})
				endTime := time.Now()
				if err != nil {
					// reader has failed, most likely due to a failed replica, we therefore do not set client to ready again
					return
				}
				r.results = append(r.results,
					Result{
						Latency:   endTime.Sub(startTime).Microseconds(),
						Timestamp: endTime.UnixMilli(),
						ID:        r.nodeID,
					})

			case <-ctx.Done():
				break loop
			}
		}
	}()
}

func (r *Reader) Send() bool {
	select {
	case r.pending <- struct{}{}:
		return true
	default:
		return false
	}
}

func (r *Reader) GetResults() []Result {
	return r.results
}
