package consts

import "time"

const (
	GorumsDialTimeout       = 5 * time.Second
	RPCContextTimeout       = 5 * time.Second
	HeartbeatSendInterval   = 1 * time.Second
	FailureDetectorInterval = 5 * time.Second
)
