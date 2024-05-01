package consts

import (
	"log/slog"
	"time"
)

const (
	GorumsDialTimeout        = 3 * time.Second
	RPCContextTimeout        = 5 * time.Second
	HeartbeatSendInterval    = 1 * time.Second
	FailureDetectionInterval = 5 * HeartbeatSendInterval
	Fanout                   = 2
	TimeFormat               = "[11:11:11]"
	LogLevel                 = slog.LevelInfo
	LogFolder                = "logs"
	NumEventHandlers         = 2
	NumTaskHandlers          = 2
	CloseMgrDelay            = 2 * RPCContextTimeout // Time to wait before disposing manager and its connections
	EventBusQueueLength      = 256
	StartupDelay             = 3 * RPCContextTimeout
)
