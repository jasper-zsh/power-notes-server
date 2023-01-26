package logic

import (
	"context"

	"powernotes-server/gateway/internal/svc"
	"powernotes-server/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SwapFlowNoteRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSwapFlowNoteRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SwapFlowNoteRelationLogic {
	return &SwapFlowNoteRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SwapFlowNoteRelationLogic) SwapFlowNoteRelation(req *types.SwapFlowNoteRequest) error {
	err := svc.SwapFlowNoteRelation(req)
	if err != nil {
		return err
	}
	return nil
}
