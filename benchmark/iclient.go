package benchmark

import "context"

type IClient interface {
	GetResults() []Result
	Send() bool
	Run(ctx context.Context)
}
