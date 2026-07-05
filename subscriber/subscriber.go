package subscriber

import (
	"forum/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MQSubHandler struct {
	svcCtx *svc.ServiceContext
	// bk     broker.Broker
}

func NewMQSubHandler(svcCtx *svc.ServiceContext) *MQSubHandler {
	// kafka 初始化
	return &MQSubHandler{
		svcCtx: svcCtx,
		// bk:     svcCtx.Kafka,
	}
}

func (s *MQSubHandler) Subscriber() {
	logx.Info("forum subscriber started")
}
