package storage

type PerNodeGossip struct {
	NodeID         string
	Key            int64
	AggValue       int64
	AggTimestamp   int64
	LocalValue     int64
	LocalTimestamp int64
}
