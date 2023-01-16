package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/neffos"
	"github.com/sirupsen/logrus"
)

type Broadcaster struct {
	connMap     map[string]*neffos.NSConn
	keyConnsMap map[string]map[string]struct{}
	connKeysMap map[string]map[string]struct{}
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		connMap:     make(map[string]*neffos.NSConn),
		keyConnsMap: make(map[string]map[string]struct{}),
		connKeysMap: make(map[string]map[string]struct{}),
	}
}

func (b *Broadcaster) Subscribe(conn *neffos.NSConn, key string) {
	connID := conn.String()
	b.connMap[connID] = conn
	keys, ok := b.connKeysMap[connID]
	if !ok {
		keys = make(map[string]struct{})
		b.connKeysMap[connID] = keys
	}
	keys[key] = struct{}{}
	connIDs, ok := b.keyConnsMap[key]
	if !ok {
		connIDs = make(map[string]struct{})
		b.keyConnsMap[key] = connIDs
	}
	connIDs[connID] = struct{}{}
}

func (b *Broadcaster) Unsubscribe(conn *neffos.NSConn, key string) {
	connID := conn.String()
	delete(b.connKeysMap[connID], key)
	delete(b.keyConnsMap[key], connID)
	if len(b.keyConnsMap[key]) == 0 {
		delete(b.keyConnsMap, key)
	}
}

func (b *Broadcaster) Disconnect(connID string) {
	delete(b.connMap, connID)
	keys, ok := b.connKeysMap[connID]
	if ok {
		for key, _ := range keys {
			connIDs, ok := b.keyConnsMap[key]
			if ok {
				delete(connIDs, connID)
			}
		}
		delete(b.connKeysMap, connID)
	}
}

func (b *Broadcaster) Broadcast(key string, event string, payload interface{}, excepts ...*neffos.NSConn) error {
	connIDs, ok := b.keyConnsMap[key]
	if !ok {
		return nil
	}
	exceptsMap := make(map[string]struct{})
	for _, except := range excepts {
		exceptsMap[except.String()] = struct{}{}
	}
	for connID, _ := range connIDs {
		if _, ok := exceptsMap[connID]; ok {
			continue
		}
		conn, ok := b.connMap[connID]
		if !ok {
			logrus.Warnf("Conn %s not found", connID)
			continue
		}
		if ok {
			err := PushToClient(conn, event, payload)
			if err != nil {
				fmt.Printf("WARN: Failed to push to %s\n", connID)
				continue
			}
		}
	}
	return nil
}

func PushToClient(conn *neffos.NSConn, event string, payload interface{}) error {
	raw := map[string]interface{}{
		"event":   event,
		"payload": payload,
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	err = conn.Conn.Socket().WriteText(b, 0)
	if err != nil {
		return err
	}
	logrus.Infof("Pushed to %s: %s", conn.String(), b)
	return nil
}
