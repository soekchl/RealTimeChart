package controllers

import (
	"net/http"
	"sync"

	. "github.com/soekchl/myUtils"
	"github.com/soekchl/websocket"
)

var (
	websocketList      []*websocket.Conn
	websocketListMutex sync.RWMutex
	wsPort             = ":8081"
)

func init() {
	http.Handle("/ws", websocket.Handler(websocketPorcess))

	go func() {
		Notice("监听 ", wsPort, " 端口.")
		err := http.ListenAndServe(wsPort, nil)
		if err != nil {
			Error(err)
		}
	}()
}

func websocketPorcess(ws *websocket.Conn) {
	websocketListMutex.Lock()
	websocketList = append(websocketList, ws)
	websocketListMutex.Unlock()

	defer func() {
		Notice("断开连接.")
		defer ws.Close()
	}()

	var (
		buff []byte
		err  error
	)

	for {
		err = websocket.Message.Receive(ws, &buff)
		if err != nil {
			Error(err)
			break
		}
		Debug("接受数据：", buff)
	}
}

func sendWebSocket(buff []byte) {
	websocketListMutex.RLock()
	for _, v := range websocketList {
		websocket.Message.Send(v, buff)
	}
	websocketListMutex.RUnlock()
}
