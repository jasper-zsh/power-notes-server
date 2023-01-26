package websocket

import (
	"github.com/gorilla/websocket"
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
	"strconv"
	"time"
)

const (
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebsocketHandler func(conn *WSConn, message []byte)
type OnDisconnectFunc func(conn *WSConn)

func RegisterHandlers(server *rest.Server) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method: http.MethodGet,
				Path:   "/ws",
				Handler: func(writer http.ResponseWriter, request *http.Request) {
					conn, err := upgrader.Upgrade(writer, request, nil)
					if err != nil {
						logx.Errorf("Failed to upgrade websocket connection %v", err)
						return
					}

					wsConn := newWSConn(conn)
					wsConn.handler = OnMessage
					OnConnect(wsConn)
					wsConn.Start()
				},
			},
		},
	)
}

type WSConn struct {
	ID           string
	ticker       *time.Ticker
	conn         *websocket.Conn
	handler      WebsocketHandler
	writeChan    chan []byte
	onDisconnect []OnDisconnectFunc
}

func newWSConn(conn *websocket.Conn) *WSConn {
	return &WSConn{
		ID:           genID(),
		ticker:       time.NewTicker(pingPeriod),
		conn:         conn,
		writeChan:    make(chan []byte, 512),
		onDisconnect: make([]OnDisconnectFunc, 0),
	}
}

func (c *WSConn) OnDisconnect(f OnDisconnectFunc) {
	c.onDisconnect = append(c.onDisconnect, f)
}

func (c *WSConn) Start() {
	go c.readPump()
	go c.writePump()
}

func (c *WSConn) Stop() {
	c.ticker.Stop()
	for _, f := range c.onDisconnect {
		f(c)
	}
	_ = c.conn.Close()
}

func (c *WSConn) Write(message []byte) {
	c.writeChan <- message
}

func (c *WSConn) readPump() {
	defer func() {
		c.Stop()
	}()
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			logx.Errorf("Failed to read message from websocket connection %v", err)
			return
		}

		if c.handler != nil {
			go c.handler(c, p)
		}
	}
}

func (c *WSConn) writePump() {
	defer func() {
		c.Stop()
	}()
	for {
		select {
		case <-c.ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case msg := <-c.writeChan:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logx.Errorf("Failed to get next writer %v", err)
				return
			}
			_, err = w.Write(msg)
			if err != nil {
				logx.Errorf("Failed to write to websocket writer %v", err)
				return
			}

			if err := w.Close(); err != nil {
				logx.Errorf("Failed to close websocket writer %v", err)
				return
			}
		}
	}
}

func genID() string {
	id, err := uuid.NewV4()
	if err != nil {
		return strconv.FormatInt(time.Now().Unix(), 10)
	}
	return id.String()
}
