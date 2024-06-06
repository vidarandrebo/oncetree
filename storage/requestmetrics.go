package storage

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/vidarandrebo/oncetree/consts"
)

type RequestMetrics struct {
	Write         atomic.Int64
	Read          atomic.Int64
	Accept        atomic.Int64
	Gossip        atomic.Int64
	Prepare       atomic.Int64
	GossipDropped atomic.Int64
	GossipSent    atomic.Int64
	LastTs        time.Time
}

type RequestRates struct {
	WritesPerSecond         float64
	ReadsPerSecond          float64
	AcceptsPerSecond        float64
	GossipsPerSecond        float64
	PreparesPerSecond       float64
	GossipsSentPerSecond    float64
	GossipsDroppedPerSecond float64

	Timestamp int64
}

func NewRequestMetrics() *RequestMetrics {
	rm := RequestMetrics{}
	return &rm
}

func (rm *RequestMetrics) Run(id string) {
	if flag.Lookup("test.v") != nil {
		return
	}
	fileName := filepath.Join(consts.LogFolder, fmt.Sprintf("replica_%s_request_rates.csv", id))
	currDir, _ := os.Getwd()
	file, err := os.Create(fileName)
	if err != nil {
		panic("file create error" + err.Error() + currDir)
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	headers := []string{"Writes", "Reads", "Accepts", "Gossips", "Prepares", "GossipsSent", "GossipsDropped", "Timestamp"}
	csvWriter.Write(headers)
	csvWriter.Flush()

	for {
		select {
		case <-time.After(time.Second):
			requestRates := rm.CollectRequestRates()
			rm.WriteRequestRates(csvWriter, requestRates)
		}
	}
}

func (rm *RequestMetrics) WriteRequestRates(writer *csv.Writer, rates RequestRates) {
	err := writer.Write([]string{
		fmt.Sprintf("%f", rates.WritesPerSecond),
		fmt.Sprintf("%f", rates.ReadsPerSecond),
		fmt.Sprintf("%f", rates.AcceptsPerSecond),
		fmt.Sprintf("%f", rates.GossipsPerSecond),
		fmt.Sprintf("%f", rates.PreparesPerSecond),
		fmt.Sprintf("%f", rates.GossipsSentPerSecond),
		fmt.Sprintf("%f", rates.GossipsDroppedPerSecond),
		fmt.Sprintf("%d", rates.Timestamp),
	})
	if err != nil {
		panic("file write error" + err.Error())
	}
	writer.Flush()
}

func (rm *RequestMetrics) CollectRequestRates() RequestRates {
	write := rm.Write.Swap(0)
	read := rm.Read.Swap(0)
	accept := rm.Accept.Swap(0)
	gossip := rm.Gossip.Swap(0)
	prepare := rm.Prepare.Swap(0)
	gossipSent := rm.GossipSent.Swap(0)
	gossipDropped := rm.GossipDropped.Swap(0)
	lastTs := rm.LastTs
	timeNow := time.Now()
	rm.LastTs = timeNow
	timeStep := timeNow.Sub(lastTs).Seconds()
	rr := RequestRates{
		WritesPerSecond:         float64(write) / timeStep,
		ReadsPerSecond:          float64(read) / timeStep,
		AcceptsPerSecond:        float64(accept) / timeStep,
		GossipsPerSecond:        float64(gossip) / timeStep,
		PreparesPerSecond:       float64(prepare) / timeStep,
		GossipsSentPerSecond:    float64(gossipSent) / timeStep,
		GossipsDroppedPerSecond: float64(gossipDropped) / timeStep,
		Timestamp:               timeNow.UnixMilli(),
	}
	return rr
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

func (rm *RequestMetrics) CountGossipDropped() {
	rm.GossipDropped.Add(1)
}

func (rm *RequestMetrics) CountGossipSent() {
	rm.GossipSent.Add(1)
}
