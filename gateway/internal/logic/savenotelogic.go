package logic

import (
	"context"

	"powernotes-server/gateway/internal/svc"
	"powernotes-server/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveNoteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveNoteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveNoteLogic {
	return &SaveNoteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveNoteLogic) SaveNote(req *types.Note) (resp *types.Note, err error) {
	note, err := svc.SaveNote(req.ToModel())
	if err != nil {
		return nil, err
	}
	return types.NewNoteFromModel(note), nil
}
