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
	TimeFormat             = "[11:11:11]"
	LogLevel               = slog.LevelInfo
	LogFolder              = "logs"
	NumEventHandlers       = 4
	NumTaskHandlers        = 16
	CloseMgrDelay          = 2 * RPCContextTimeout // Time to wait before disposing manager and its connections
	EventBusQueueLength    = 2
	StartupDelay           = 1 * RPCContextTimeout
	TestWaitAfterWrite     = 10 * time.Second
)
