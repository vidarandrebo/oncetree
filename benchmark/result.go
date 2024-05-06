package benchmark

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vidarandrebo/oncetree/consts"
)

type Result struct {
	Latency   int64
	Timestamp int64
	ID        string
	Type      OperationType
}

func WriteResultsToDisk(results []Result, id string) {
	writeFilename := filepath.Join(consts.LogFolder, fmt.Sprintf("client_%s_write_results.csv", id))
	readFilename := filepath.Join(consts.LogFolder, fmt.Sprintf("client_%s_read_results.csv", id))
	writesFile, err := os.Create(writeFilename)
	if err != nil {
		panic("create file error")
	}
	readsFile, err := os.Create(readFilename)
	if err != nil {
		panic("create file error")
	}
	defer writesFile.Close()
	defer readsFile.Close()
	csvWrites := csv.NewWriter(writesFile)
	csvReads := csv.NewWriter(readsFile)
	defer csvWrites.Flush()
	defer csvReads.Flush()
	headers := []string{"Timestamp", "NodeID", "Latency"}
	err = csvWrites.Write(headers)
	err = csvReads.Write(headers)
	for _, result := range results {
		if result.Type == Read {
			err = csvReads.Write([]string{
				fmt.Sprintf("%d", result.Timestamp),
				result.ID,
				fmt.Sprintf("%d", result.Latency),
			})
		} else if result.Type == Write {
			err = csvWrites.Write([]string{
				fmt.Sprintf("%d", result.Timestamp),
				result.ID,
				fmt.Sprintf("%d", result.Latency),
			})
		}
		if err != nil {
			panic("write error: " + err.Error())
		}
	}
}
