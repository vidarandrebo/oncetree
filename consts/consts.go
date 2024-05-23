package consts

import (
	"log/slog"
	"time"
)

const (
	GorumsDialTimeout      = 3 * time.Second
	RPCContextTimeout      = 5 * time.Second
	HeartbeatSendInterval  = 1 * time.Second
	FailureDetectorStrikes = 5
	Fanout                 = 2
	LogLevel               = slog.LevelInfo
	LogFolder              = "logs"
	NumEventHandlers       = 4
	NumTaskHandlers        = 4
	CloseMgrDelay          = 3 * RPCContextTimeout // Time to wait before disposing manager and its connections
	EventBusQueueLength    = 256
	StartupDelay           = 1 * RPCContextTimeout
	TestWaitAfterWrite     = 5 * time.Second
	GossipWorkerBuffSize   = 10000
	BenchmarkTime          = 60 * time.Second
	BenchmarkNumKeys       = 5000
)
