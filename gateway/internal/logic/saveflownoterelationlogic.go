package logic

import (
	"context"

	"powernotes-server/gateway/internal/svc"
	"powernotes-server/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveFlowNoteRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveFlowNoteRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveFlowNoteRelationLogic {
	return &SaveFlowNoteRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveFlowNoteRelationLogic) SaveFlowNoteRelation(req *types.FlowNoteRelation) (resp *types.FlowNoteRelation, err error) {
	rel, err := svc.SaveFlowNoteRelation(req.ToModel())
	if err != nil {
		return nil, err
	}
	return types.NewFlowNoteRelationFromModel(rel), nil
}
