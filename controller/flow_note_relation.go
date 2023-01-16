package controller

import (
	"github.com/kataras/iris/v12/context"
	"powernotes-server/model"
	"powernotes-server/service"
)

func SaveFlowNoteRelation(ctx *context.Context) {
	rel := &model.FlowNoteRelation{}
	err := ctx.ReadJSON(rel)
	if err != nil {
		RenderJSON(ctx, nil, err)
		return
	}
	rel, err = service.SaveFlowNoteRelation(rel)
	RenderJSON(ctx, rel, err)
}

func RemoveFlowNoteRelation(ctx *context.Context) {
	params := ctx.Params()
	flowID, _ := params.GetInt64("flow_id")
	noteID, _ := params.GetInt64("note_id")
	err := service.RemoveFlowNoteRelation(flowID, noteID)
	RenderJSON(ctx, nil, err)
}

func SwapFlowNoteRelation(ctx *context.Context) {
	req := &service.SwapFlowNote{}
	err := ctx.ReadJSON(req)
	if err != nil {
		RenderJSON(ctx, nil, err)
		return
	}
	err = service.SwapFlowNoteRelation(req)
	RenderJSON(ctx, nil, err)
}
