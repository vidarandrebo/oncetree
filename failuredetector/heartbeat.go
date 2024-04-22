package failuredetector

import (
	"log/slog"

	"github.com/relab/gorums"
	fdprotos "github.com/vidarandrebo/oncetree/protos/failuredetector"
)

func (fd *FailureDetector) Heartbeat(ctx gorums.ServerCtx, request *fdprotos.HeartbeatMessage) {
	fd.logger.Debug(
		"RPC Heartbeat",
		slog.String("id", request.GetNodeID()),
	)
	fd.alive.Increment(request.GetNodeID(), 1)
	if fd.suspected.Contains(request.GetNodeID()) {
		fd.logger.Error("received heartbeat from suspected node", slog.String("id", request.GetNodeID()))
		panic("heartbeat problem")
	}
}
