package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"powernotes-server/config"
	"powernotes-server/controller"
	"powernotes-server/handler"
	"powernotes-server/model"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	model.InitDB(cfg.DB)
	app := iris.New()

	ws := neffos.New(websocket.DefaultGorillaUpgrader, neffos.Events{
		neffos.OnNativeMessage: handler.OnNativeMessage,
	})
	ws.OnDisconnect = handler.OnDisconnect
	app.Get("/ws", websocket.Handler(ws))

	app.Post("/flow", controller.SaveFlow)
	app.Delete("/flow/{id}", controller.RemoveFlow)
	app.Post("/note", controller.SaveNote)
	app.Delete("/note/{id}", controller.RemoveNote)
	app.Post("/flow_note_relation", controller.SaveFlowNoteRelation)
	app.Delete("/flow/{flow_id}/note_relation/{note_id}", controller.RemoveFlowNoteRelation)
	app.Post("/flow_note_relation/swap", controller.SwapFlowNoteRelation)

	err = app.Listen(":5555")
	if err != nil {
		panic(err)
	}
}
