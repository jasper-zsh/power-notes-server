package controller

import (
	"github.com/kataras/iris/v12/context"
	"gorm.io/gorm/clause"
	"powernotes-server/handler"
	"powernotes-server/model"
	"powernotes-server/service"
)

func SaveFlow(ctx *context.Context) {
	flow := model.Flow{}
	err := ctx.ReadJSON(&flow)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	result := model.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Save(&flow)
	if result.Error != nil {
		ctx.StopWithError(500, result.Error)
		return
	}
	handler.ProjectBroadcaster.Broadcast(flow.ProjectName, "flow", &flow)
	err = ctx.JSON(&flow)
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
}

func RemoveFlow(ctx *context.Context) {
	params := ctx.Params()
	id, _ := params.GetInt64("id")
	flow, err := service.RemoveFlow(id)
	RenderJSON(ctx, flow, err)
}
