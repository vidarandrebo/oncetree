package storage

import "fmt"

type TimestampedValue struct {
	Value     int64
	Timestamp int64
}

func (tsv TimestampedValue) String() string {
	return fmt.Sprintf("(Value: %d, ts: %d)", tsv.Value, tsv.Timestamp)
}
