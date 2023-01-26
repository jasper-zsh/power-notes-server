package websocket

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
)

type Handler func(conn *WSConn, body map[string]interface{}) error

type Event struct {
	Event   string                 `json:"event"`
	Payload map[string]interface{} `json:"payload"`
}

var handlers = map[string]Handler{
	"open_file":     OnOpenFile,
	"close_file":    OnCloseFile,
	"open_project":  OnOpenProject,
	"close_project": OnCloseProject,
}

func OnConnect(c *WSConn) {
	c.OnDisconnect(func(conn *WSConn) {
		FileBroadcaster.Disconnect(conn.ID)
		ProjectBroadcaster.Disconnect(conn.ID)
	})
}

func OnMessage(c *WSConn, msg []byte) {
	logx.Infof("Conn: %s   Msg: %s", c.ID, msg)
	event := Event{}
	err := json.Unmarshal(msg, &event)
	if err != nil {
		logx.Errorf("Failed to unmarshal json %v", err)
		return
	}
	handler, ok := handlers[event.Event]
	if ok {
		err = handler(c, event.Payload)
		if err != nil {
			logx.Errorf("Failed to run handler %s %v", event.Event, err)
			return
		}
	} else {
		logx.Errorf("Event %s undefined.", event.Event)
	}
}
