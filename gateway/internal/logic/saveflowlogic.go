package logic

import (
	"context"

	"powernotes-server/gateway/internal/svc"
	"powernotes-server/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveFlowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveFlowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveFlowLogic {
	return &SaveFlowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveFlowLogic) SaveFlow(req *types.Flow) (resp *types.Flow, err error) {
	flow, err := svc.SaveFlow(req.ToModel())
	if err != nil {
		return nil, err
	}
	return types.NewFlowFromModel(flow), nil
}
