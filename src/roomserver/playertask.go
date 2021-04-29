package main

import (
	"base/gonet"
	"github.com/gorilla/websocket"
	"time"
)

/*************************
PlayerTask 并通过PlayTaskMgr 管理PlayerTask
将 *websocket.Conn 转化为 wstask *gonet.WebSocketTask
并封装成PlayerTask
*************************/

type PlayerTask struct {
	wstask *gonet.WebSocketTask //用户的websocket链接
	id     uint32               //在 初始化//Todo
	//todo 创建房间与场景
	//room       *Room                //所属房间
	//scene      *Scene               //玩家场景
	activetime time.Time //活跃时间
	angle      uint32    //角度

}

func NewPlayerTask(conn *websocket.Conn) *PlayerTask {
	temPTask := &PlayerTask{
		wstask:     gonet.NewWebSocketTask(conn),
		activetime: time.Now(), //开始时间
		//todo 初始化场景
	}

	// 需要实现 ParseMsg() OnClose()
	temPTask.wstask.Derived = temPTask

	return temPTask
}

func (PT *PlayerTask) Start() {

}

func (PT *PlayerTask) ParseMsg(data []byte, flag byte) bool {
	panic("implement me")
}

func (PT *PlayerTask) OnClose() {
	panic("implement me")
}
