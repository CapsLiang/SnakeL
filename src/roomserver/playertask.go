package main

import (
	"base/gonet"
	"common"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"math/rand"
	"runtime/debug"
	"sync"
	"time"
)

/*************************
PlayerTask 并通过PlayTaskMgr 管理PlayerTask
将 *websocket.Conn 转化为 wstask *gonet.WebSocketTask
并封装成PlayerTask
*************************/

type PlayerTask struct {
	wstask *gonet.WebSocketTask //用户的websocket链接
	id     uint32               //在Start()中初始化
	//todo 创建房间与场景
	room       *Room     //所属房间
	scene      *Scene    //玩家场景
	activetime time.Time //活跃时间
	angle      uint32    //角度 todo 在parsemsg中获得

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
	PT.id = rand.New(rand.NewSource(time.Now().UnixNano())).Uint32() % 100 // 待优化
	PT.wstask.Start()                                                      //开两个goroutine收发消息
	PT.wstask.Verify()                                                     // 使用验证通常是为了防止用户连接而不使用，占据服务器资源 验证客户端是否合法后可以减少这种情况

	PlayerTaskMgr_GetMe().Add(PT) //添加到PTmgr管理
	//todo 分配房间

	//this.scene.self.Id = this.id
	//room, err := RoomMgr_GetMe().GetRoom(PT)
	//if nil != err {
	//	glog.Error("[roomserver] Allocate room fail ", err)
	//	return
	//}
	//PT.scene.room = room

}

//todo ParseMsg
func (this *PlayerTask) ParseMsg(data []byte, flag byte) bool {
	glog.Info("[WS] Parse Msg", data)
	this.activetime = time.Now()

	//todo msg需要约定
	msgtype := common.MsgType(uint16(data[2]))

	switch msgtype {

	}

	return true
}

//todo OnClose()
func (PT *PlayerTask) OnClose() {

	PT.wstask.Close()
	PlayerTaskMgr_GetMe().Del(PT)

	PT.room = nil
	PT.scene = nil
}

func (this *PlayerTask) SendSceneMsg() bool {
	if nil == this.scene {
		return false
	}

	msg := this.scene.SceneMsg()
	if nil == msg {
		glog.Error("[Scene] Msg nil")
		return false
	}
	//todo 传输现在为json
	return this.wstask.AsyncSend(msg, 0)
}

/*************************
通过PlayTaskMgr 管理PlayerTask
*************************/
/*
//开一个协程管理PlayerTask
func PlayerTaskMgr_GetMe() *PlayerTaskMgr
//超时断开连接
func (thisPTMgr *PlayerTaskMgr) iTimeAction()
//删除管理并断开连接
func (thisPTMgr *PlayerTaskMgr) Del(PTask *PlayerTask) bool
//添加到map管理
func (thisPTMgr *PlayerTaskMgr) Add(PTask *PlayerTask) bool
//get id所对应的连接
func (thisPTMgr *PlayerTaskMgr) Get(id uint32) *PlayerTask
*/

type PlayerTaskMgr struct {
	mutex sync.RWMutex //读写锁
	tasks map[uint32]*PlayerTask
}

var mPlayerTaskMgr *PlayerTaskMgr

//开一个协程管理PlayerTask
func PlayerTaskMgr_GetMe() *PlayerTaskMgr {
	//初始化 如果没有mPlayerTaskMgr则创建 有则直接返回
	if nil == mPlayerTaskMgr {
		mPlayerTaskMgr = &PlayerTaskMgr{
			//通过id获取PlayerTask
			tasks: make(map[uint32]*PlayerTask),
		}
		go mPlayerTaskMgr.iTimeAction() //开始管理连接
	}

	return mPlayerTaskMgr
}

func (thisPTMgr *PlayerTaskMgr) iTimeAction() {
	var (
		timeTicker = time.NewTicker(time.Second)
		loop       uint64
		ptasks     []*PlayerTask
	)

	defer func() {
		timeTicker.Stop()
		if err := recover(); nil != err {
			glog.Error("[异常] 定时线程出错", err, "\n", string(debug.Stack()))
		}
	}()

	for {
		select {
		case <-timeTicker.C:
			if 0 == loop%5 {
				now := time.Now()
				thisPTMgr.mutex.RLock()
				for _, t := range thisPTMgr.tasks {
					if now.Sub(t.activetime) > common.Task_TimeOut*time.Second {
						//超时链接
						ptasks = append(ptasks, t)
					}
				}
				thisPTMgr.mutex.RUnlock()

				for _, t := range ptasks {
					if !t.wstask.Stop() {
						thisPTMgr.Del(t) //删除超时链接
					}
					glog.Info("[Player] Connection timeout, player id=", t.id) //连接超时
				}
				ptasks = ptasks[:0] //置空
			}
			loop++
		}
	}

}

//删除管理并断开连接
func (thisPTMgr *PlayerTaskMgr) Del(PTask *PlayerTask) bool {
	if nil == PTask {
		glog.Error("[WS] Player Task Manager Del Fail,Player Task is Nil")
		return false
	}

	thisPTMgr.mutex.Lock()
	defer thisPTMgr.mutex.Unlock()

	//根据id取PlayerTask
	task, ok := thisPTMgr.tasks[PTask.id]
	if !ok {
		return false
	}
	if PTask != task {
		glog.Error("[WS] Player Task Manager Del Fail, ", PTask.id, ",", &PTask, ",", &task)
		return false
	}
	delete(thisPTMgr.tasks, PTask.id)
	return true
}

//添加到map管理
func (thisPTMgr *PlayerTaskMgr) Add(PTask *PlayerTask) bool {
	if nil == PTask {
		glog.Error("[WS] Player Task Manager Add Fail, Nil")
		return false
	}

	thisPTMgr.mutex.Lock()
	defer thisPTMgr.mutex.Unlock()

	thisPTMgr.tasks[PTask.id] = PTask

	return true
}

//get id所对应的连接
func (thisPTMgr *PlayerTaskMgr) Get(id uint32) *PlayerTask {
	thisPTMgr.mutex.RLock()
	defer thisPTMgr.mutex.RUnlock()

	PTask, ok := thisPTMgr.tasks[id]
	if !ok {
		return nil
	}

	return PTask
}
