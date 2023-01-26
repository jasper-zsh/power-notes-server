package logic

import (
	"context"

	"powernotes-server/gateway/internal/svc"
	"powernotes-server/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveFlowLogic {
	return &RemoveFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveFlowLogic) RemoveFlow(req *types.IDRequest) (resp *types.Flow, err error) {
	flow, err := svc.RemoveFlow(req.ID)
	if err != nil {
		return nil, err
	}
	return types.NewFlowFromModel(flow), nil
}
