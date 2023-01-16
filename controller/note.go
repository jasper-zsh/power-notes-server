package controller

import (
	"github.com/kataras/iris/v12/context"
	"powernotes-server/model"
	"powernotes-server/service"
)

func SaveNote(ctx *context.Context) {
	note := &model.Note{}
	err := ctx.ReadJSON(note)
	if err != nil {
		RenderJSON(ctx, nil, err)
		return
	}
	note, err = service.SaveNote(note)
	RenderJSON(ctx, note, err)
}

func RemoveNote(ctx *context.Context) {
	params := ctx.Params()
	id, _ := params.GetInt64("id")
	note, err := service.RemoveNote(id)
	RenderJSON(ctx, note, err)
}
