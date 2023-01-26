package logic

import (
	"context"

	"powernotes-server/gateway/internal/svc"
	"powernotes-server/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveNoteLogic {
	return &RemoveNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveNoteLogic) RemoveNote(req *types.IDRequest) (resp *types.Note, err error) {
	note, err := svc.RemoveNote(req.ID)
	if err != nil {
		return nil, err
	}
	return types.NewNoteFromModel(note), nil
}
