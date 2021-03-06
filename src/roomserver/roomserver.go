package main

import (
	"base/env"
	"base/gonet"
	"flag"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
)

type RoomServer struct {
	gonet.Service
	roomser *gonet.WebSocketServer
}

var mServer *RoomServer

func RoomServer_GetMe() *RoomServer {
	if nil == mServer {
		mServer = &RoomServer{
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
	glog.Info("[WS] Connected websocket已连接")
	//创建玩家连接 分配房间
	NewPlayerTask(conn).Start()

}

func (r RoomServer) Init() bool {
	//RoomGrpcClient_GetMe().Init() 通过grpc获取此房间的ip port
	if !RoomGrpcClient_GetMe().Init() {
		glog.Error("[gRPC] 房间服务客户端初始化失败")
		return false
	}

	go func() {
		//开一条goroutine 通过gonet提供的接口读取配置并监听
		err := r.roomser.WSBind(env.Get("room", "listen"))
		if nil != err {
			glog.Error("[Start] Bind Port 失败")
		}
	}()

	return true
}

func (r RoomServer) Reload() {
}

func (r RoomServer) MainLoop() {
	//可以不用
	time.Sleep(time.Second)
}

func (r RoomServer) Final() bool {
	//关闭RPCClient
	RoomGrpcClient_GetMe().Close()
	return true
}

//执行时通过命令行初始化配置环境
var (
	config  = flag.String("config", "", "config file")
	logfile = flag.String("logfile", "", "log file")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().Unix())
	env.Load(*config)

	if "" != *logfile {
		glog.SetLogFile(*logfile)
	} else {
		glog.SetLogFile(env.Get("room", "log"))
	}
	defer glog.Flush()

	RoomServer_GetMe().Main()
}
