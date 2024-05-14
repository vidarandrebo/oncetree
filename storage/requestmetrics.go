package storage

import "sync/atomic"

type RequestMetrics struct {
	Write   atomic.Int64
	Read    atomic.Int64
	Accept  atomic.Int64
	Gossip  atomic.Int64
	Prepare atomic.Int64
}

func NewRequestMetrics() *RequestMetrics {
	rm := RequestMetrics{}
	return &rm
}

func (rm *RequestMetrics) Reset() {
}

func (rm *RequestMetrics) CountWrite() {
	rm.Write.Add(1)
}

func (rm *RequestMetrics) CountRead() {
	rm.Read.Add(1)
}

func (rm *RequestMetrics) CountAccept() {
	rm.Accept.Add(1)
}

func (rm *RequestMetrics) CountGossip() {
	rm.Gossip.Add(1)
}

func (rm *RequestMetrics) CountPrepare() {
	rm.Prepare.Add(1)
}
