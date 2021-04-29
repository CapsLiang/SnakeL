package roomserver

import (
	"base/gonet"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type RoomServer struct {
	gonet.Service
	roomser *gonet.WebSocketServer
}

var mServer *RoomServer

func RoomServer_GetMe() *RoomServer {
	if nil == mServer {
		mServer = &RoomServer{
			Service: gonet.Service{},
			roomser: &gonet.WebSocketServer{},
		}
	}

	// 需要实现  Init() Reload() MainLoop() Final() 接口
	mServer.Derived = mServer

	//需要实现OnWSAccept方法
	mServer.roomser.Derived = mServer

	return mServer
}

func (r RoomServer) OnWSAccept(conn *websocket.Conn) {
	//Todo: NewPlayerTask(conn).Start()
	glog.Info("[WS] Connected websocket已连接")
}
