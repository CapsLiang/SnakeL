package main

import (
	"base/env"
	"errors"
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
	id          uint32                 //房间id newroom中初始化
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
		id:       rid,
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
	this.GameLoop() //开始游戏主循环
}

func (this *Room) GameLoop() {
	timeTicker := time.NewTicker(time.Millisecond * 10)
	stop := false
	for !stop {
		// SceneMsg, 用于同步场景
		select {
		case <-timeTicker.C:
			if this.timeloop%2 == 0 { //0.02s
				//todo room update()
				this.update()
				//todo room sendRoomMsg()
				this.sendRoomMsg()
			}

			if this.timeloop%100 == 0 { //1s
				//todo sendTime()
				this.sendTime(this.totgametime - this.timeloop/100)
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
		glog.Info("[Room] ", player.id, "玩家已经在[", this.id, "]房间里面了")
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
		p.SendSceneMsg()
	}
}

func (this *Room) sendTime(t uint64) {

}

func (this *Room) update() {

}
