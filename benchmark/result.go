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
}

func WriteResultsToDisk(results []Result, id string, operationType OperationType) {
	fileName := ""
	if operationType == Read {
		fileName = filepath.Join(consts.LogFolder, fmt.Sprintf("client_%s_read_results.csv", id))
	} else if operationType == Write {
		fileName = filepath.Join(consts.LogFolder, fmt.Sprintf("client_%s_write_results.csv", id))
	}
	file, err := os.Create(fileName)
	if err != nil {
		panic("create file error")
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()
	headers := []string{"Timestamp", "NodeID", "Latency"}
	err = csvWriter.Write(headers)
	for _, result := range results {
		err = csvWriter.Write([]string{
			fmt.Sprintf("%d", result.Timestamp),
			result.ID,
			fmt.Sprintf("%d", result.Latency),
		})
		if err != nil {
			panic("write error: " + err.Error())
		}
	}
}
