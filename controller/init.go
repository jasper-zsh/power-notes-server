package controller

import "github.com/kataras/iris/v12/context"

func RenderJSON(ctx *context.Context, result interface{}, err error) {
	if err != nil {
		ctx.StopWithError(500, err)
		return
	}
	err = ctx.JSON(result)
	if err != nil {
		ctx.StopWithError(500, err)
	}
}
