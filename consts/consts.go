package consts

import "time"

const (
	GorumsDialTimeout       = 3 * time.Second
	RPCContextTimeout       = 3 * time.Second
	HeartbeatSendInterval   = 1 * time.Second
	FailureDetectorInterval = 5 * time.Second
)
