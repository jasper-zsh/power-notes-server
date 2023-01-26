package websocket

import (
	"encoding/json"
	"github.com/kataras/neffos"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type Broadcaster struct {
	connMap     map[string]*WSConn
	keyConnsMap map[string]map[string]struct{}
	connKeysMap map[string]map[string]struct{}
	mapsLock    sync.RWMutex
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		connMap:     make(map[string]*WSConn),
		keyConnsMap: make(map[string]map[string]struct{}),
		connKeysMap: make(map[string]map[string]struct{}),
		mapsLock:    sync.RWMutex{},
	}
}

func (b *Broadcaster) Subscribe(conn *WSConn, key string) {
	b.mapsLock.Lock()
	defer b.mapsLock.Unlock()
	connID := conn.ID
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

func (b *Broadcaster) Unsubscribe(conn *WSConn, key string) {
	b.mapsLock.Lock()
	defer b.mapsLock.Unlock()
	connID := conn.ID
	delete(b.connKeysMap[connID], key)
	delete(b.keyConnsMap[key], connID)
	if len(b.keyConnsMap[key]) == 0 {
		delete(b.keyConnsMap, key)
	}
}

func (b *Broadcaster) Disconnect(connID string) {
	b.mapsLock.Lock()
	defer b.mapsLock.Unlock()
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
	b.mapsLock.RLock()
	defer b.mapsLock.RUnlock()
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
			logx.Infof("Conn %s not found", connID)
			continue
		}
		if ok {
			err := PushToClient(conn, event, payload)
			if err != nil {
				logx.Infof("WARN: Failed to push to %s\n", connID)
				continue
			}
		}
	}
	return nil
}

func PushToClient(conn *WSConn, event string, payload interface{}) error {
	raw := map[string]interface{}{
		"event":   event,
		"payload": payload,
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	conn.Write(b)
	logx.Infof("Pushed to %s: %s", conn.ID, b)
	return nil
}
