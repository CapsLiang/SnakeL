package main

import (
	"base/env"
	"common"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"strconv"
	"time"
)

const (
	MaxPlayerNum uint32 = 2
)

//提供信息给roommgr管理房间
type Room struct {
	//mutex    sync.RWMutex
	roomid      uint32                 //房间id newroom中初始化
	roomtype    uint32                 //房间类型
	players     map[uint32]*PlayerTask //房间内的玩家
	curnum      uint32                 //当前房间内玩家数
	isstart     bool                   //此房间是否运行
	timeloop    uint64
	stopch      chan bool
	Isstop      bool
	totgametime uint64 //游戏结束时间 in second
	endchan     chan bool
}

func NewRoom(rtype, rid uint32) *Room {
	room := &Room{
		roomid:   rid,
		roomtype: rtype,
		players:  make(map[uint32]*PlayerTask),
		curnum:   0,
		isstart:  false,
		Isstop:   false,
		endchan:  make(chan bool),
	}
	room.totgametime, _ = strconv.ParseUint(env.Get("room", "time"), 10, 64)
	return room
}

func (this *Room) Start() {
	this.isstart = true

	for _, snake := range this.players {
		snake.scene.InitSnake()
	}

	fmt.Println("[roomserver] room start")
	this.GameLoop() //开始游戏主循环
}

func (this *Room) GameLoop() {
	fmt.Println("[roomserver] loop start")
	timeTicker := time.NewTicker(time.Millisecond * 10) //10ms 发一次消息
	stop := false
	for !stop {
		// SceneMsg, 用于同步场景
		select {
		case <-timeTicker.C:
			if this.timeloop%2 == 0 { //0.02s 20ms
				fmt.Println("timeloop ", this.timeloop)
				this.update()
				this.sendRoomMsg()

				fmt.Println("                                      ")
			}

			if this.timeloop%100 == 0 { //1s
				this.sendTime(this.totgametime - this.timeloop/100)
			}

			if this.timeloop%200 == 0 { //2s
				if this.curnum == 0 {
					stop = true
				}
			}

			if this.timeloop != 0 && this.timeloop%(this.totgametime*100) == 0 {
				//超时
				stop = true
			}
			this.timeloop++
			if this.Isstop {
				stop = true
			}
		}
	}
	this.Close()
}

//给一个玩家分配房间
func (this *Room) AddPlayer(player *PlayerTask) error {

	if this.checkPlayer(player) {
		glog.Info("[Room] ", player.id, "玩家已经在[", this.roomid, "]房间里面了")
		return nil
	}
	if this.curnum >= MaxPlayerNum {
		glog.Error("[Room] 房间已满")
		return errors.New("room is full")
	}

	this.curnum++
	this.players[player.id] = player
	this.players[player.id].room = this

	return nil
}

func (this *Room) IsFull() bool {
	if this.curnum < MaxPlayerNum {
		return false
	}
	return true
}

func (this *Room) Close() {

}

func (this *Room) checkPlayer(player *PlayerTask) bool {
	//不在map中return false
	if _, ok := this.players[player.id]; !ok {
		return false
	}
	return true
}

func (this *Room) sendRoomMsg() {
	for _, p := range this.players {
		sendsucced := p.SendSceneMsg()
		if sendsucced {
			fmt.Println("[sendRoomMsg] 玩家: ", p.id, " 发送消息成功")
			glog.Info("[sendRoomMsg] 玩家: ", p.id, " 发送消息成功")
		}
	}
}

//还有多久游戏结束
func (this *Room) sendTime(t uint64) {
	for _, p := range this.players {
		t := common.RetTimeMsg{
			Time: t,
		}
		jstr, err := json.Marshal(t)
		if err != nil {
			glog.Error("[Time] marshal jsonMsg err")
			return
		}
		fmt.Println(string(jstr))
		p.wstask.AsyncSend(jstr, 0)
	}
}

func (this *Room) update() {
	this.AddFoods()
	glog.Info("[room]更新食物")

	for _, p := range this.players {
		p.Update()
	}
	glog.Info("[room]更新房间")
	//for _, p := range this.players {
	//	p.UpdateOthers()
	//}

}
