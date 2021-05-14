package main

import (
	"base/gonet"
	"bytes"
	"common"
	"encoding/binary"
	"fmt"
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
	wstask     *gonet.WebSocketTask //用户的websocket链接
	name       string               //玩家输入 默认为snake
	id         uint32               //在Start()中初始化
	room       *Room                //所属房间
	scene      *Scene               //玩家场景
	activetime time.Time            //活跃时间
	direct     float64
}

func NewPlayerTask(conn *websocket.Conn) *PlayerTask {
	temPTask := &PlayerTask{
		wstask: gonet.NewWebSocketTask(conn),
		scene: &Scene{
			snake:   SnakeBody{},
			speed:   common.SceneSpeed,
			preTime: time.Now().UnixNano() / 1e6, //毫秒
		},
		activetime: time.Now(), //开始时间
	}

	// 需要实现 ParseMsg() OnClose()
	temPTask.wstask.Derived = temPTask

	return temPTask
}

func (this *PlayerTask) Start() {
	this.id = rand.New(rand.NewSource(time.Now().UnixNano())).Uint32() % 100 // 待优化
	this.name = "snake"
	this.wstask.Start()  //开两个goroutine收发消息
	this.wstask.Verify() // 使用验证通常是为了防止用户连接而不使用，占据服务器资源 验证客户端是否合法后可以减少这种情况

	PlayerTaskMgr_GetMe().Add(this) //添加到PTmgr管理

	//分配房间
	room, err := RoomMgr_GetMe().GetRoom(this)
	if nil != err {
		glog.Error("[roomserver] 申请房间失败", err)
		return
	}
	this.scene.room = room
	this.scene.snake.thisplayer = this
	//this.scene.InitSnake() //初始化蛇
}

func (this *PlayerTask) ParseMsg(data []byte, flag byte) bool {
	glog.Info("[WS] Parse Msg", data)
	this.activetime = time.Now()

	msgtype := common.MsgType(uint16(data[2]) | uint16(data[3])<<8)

	switch msgtype {
	case common.MsgType_Move:
		var angle int32
		err := binary.Read(bytes.NewReader(data[4:]), binary.LittleEndian, &angle)
		if nil != err {
			glog.Error("[WS] Endian Trans Fail", data)
			return false
		}
		glog.Info("[WS] Parse Msg Move ", angle)
		if nil == this.room {
			fmt.Println("此玩家房间未初始化")
			return false
		}
		if this.room.Isstop {
			fmt.Println("此玩家房间已经停止")
			return false
		}
		if nil == this.scene {
			fmt.Println("此玩家场景未初始化")
			return false
		}

		this.direct = float64(angle)
		this.scene.UpdateSnakePOINT(this.direct)
		this.scene.UpdateSpeed(common.SceneSpeed)

	case common.MsgType_SpeedUp:
		if nil == this.room {
			return false
		}
		if this.room.Isstop {
			return false
		}
		if nil == this.scene {
			return false
		}
		this.scene.UpdateSpeed(common.SceneSpeed * common.SceneSpeedUp)

	case common.MsgType_Finsh:
		this.room.Close()

	case common.MsgType_Heart:
		this.wstask.AsyncSend(data, flag)
	default:
	}
	return true
}

func (this *PlayerTask) Stop() bool {
	return this.wstask.Stop()
}

func (this *PlayerTask) OnClose() {
	this.wstask.Close()
	PlayerTaskMgr_GetMe().Del(this)

	this.room = nil
	this.scene = nil
}

func (this *PlayerTask) Update() {
	if nil == this.scene {
		return
	}

	this.scene.UpdateSnakePOINT(this.direct)
}

//func (this *PlayerTask) UpdateOthers() {
//	if nil == this.scene {
//		return
//	}
//
//	this.scene.UpdateOthersSnake()
//}

func (this *PlayerTask) SendSceneMsg() bool {
	if nil == this.scene {
		return false
	}

	msg := this.scene.SceneMsg()
	if nil == msg {
		glog.Error("[Scene] 消息为空")
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
					if !t.Stop() {
						thisPTMgr.Del(t) //删除超时链接
					}
					glog.Info("[Player] 连接超时, 玩家id=", t.id) //连接超时
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
		glog.Error("[roomserver][WS] Player Task Manager Del Fail,Player Task is Nil")
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
		glog.Error("[roomserver][WS] Player Task Manager Del Fail, ", PTask.id, ",", &PTask, ",", &task)
		return false
	}
	delete(thisPTMgr.tasks, PTask.id)
	return true
}

//添加到map管理
func (thisPTMgr *PlayerTaskMgr) Add(PTask *PlayerTask) bool {
	if nil == PTask {
		glog.Error("[roomserver][WS] Player Task Manager Add Fail, Nil")
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
