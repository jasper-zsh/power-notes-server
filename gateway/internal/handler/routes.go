// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"powernotes-server/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/flow",
				Handler: saveFlowHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/flow/:id",
				Handler: removeFlowHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/note",
				Handler: saveNoteHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/note/:id",
				Handler: removeNoteHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/flow_note_relation",
				Handler: saveFlowNoteRelationHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/flow/:flow_id/note_relation/:note_id",
				Handler: removeFlowNoteRelationHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/flow_note_relation/swap",
				Handler: swapFlowNoteRelationHandler(serverCtx),
			},
		},
	)
}
