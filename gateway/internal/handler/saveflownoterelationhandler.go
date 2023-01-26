package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"powernotes-server/gateway/internal/logic"
	"powernotes-server/gateway/internal/svc"
	"powernotes-server/gateway/internal/types"
)

func saveFlowNoteRelationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FlowNoteRelation
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSaveFlowNoteRelationLogic(r.Context(), svcCtx)
		resp, err := l.SaveFlowNoteRelation(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
