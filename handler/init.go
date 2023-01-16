package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/neffos"
)

type Handler func(conn *neffos.NSConn, body map[string]interface{}) error

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

func OnNativeMessage(c *neffos.NSConn, msg neffos.Message) error {
	fmt.Printf("Conn: %s   Msg: %s\n", c.String(), msg.Body)
	event := Event{}
	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		return err
	}
	handler, ok := handlers[event.Event]
	if ok {
		err = handler(c, event.Payload)
		if err != nil {
			return err
		}
	}
	return nil
}
