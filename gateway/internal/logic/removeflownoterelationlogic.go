package logic

import (
	"context"

	"powernotes-server/gateway/internal/svc"
	"powernotes-server/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveFlowNoteRelationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveFlowNoteRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveFlowNoteRelationLogic {
	return &RemoveFlowNoteRelationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveFlowNoteRelationLogic) RemoveFlowNoteRelation(req *types.FlowNoteIDRequest) error {
	err := svc.RemoveFlowNoteRelation(req.FlowID, req.NoteID)
	if err != nil {
		return err
	}
	return nil
}
