package consts

import (
	"log/slog"
	"time"
)

const (
	GorumsDialTimeout       = 3 * time.Second
	RPCContextTimeout       = 10 * time.Second
	HeartbeatSendInterval   = 1 * time.Second
	FailureDetectorInterval = 5 * time.Second
	Fanout                  = 2
	TimeFormat              = "[11:11:11]"
	LogLevel                = slog.LevelDebug
	LogFolder               = "logs"
)
