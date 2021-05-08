package main

import (
	"common"
	"encoding/json"
	"github.com/golang/glog"
	"math"
	"time"
)

type SnakeBody struct {
	id         uint32
	name       string
	thisplayer *PlayerTask
	direct     uint32
	head       common.POINT
	body       []common.POINT
	score      int32
	radius     float64
	//invincible bool todo 无敌
}

func (this *SnakeBody) SnakeDie() {

}

type FoodList struct {
	//foodMutex sync.Mutex
	foodlist []common.Food //食物列表
	eatfood  map[int]bool  //判断是否被吃 每次更新
}

//公用食物
var mFoods *FoodList

type Scene struct {
	room *Room //所属房间

	snake SnakeBody //本条蛇

	others []SnakeBody //其他人的信息

	preTime int64 //上一帧时间
}

func (this *Scene) AddFoods() {
	//不存在时 创建食物数组
	if nil == mFoods {
		mFoods = &FoodList{
			//foodMutex: sync.Mutex{},
			foodlist: make([]common.Food, 0, common.FoodNum),
			eatfood:  make(map[int]bool),
		}
		//初始化食物数组
		for i := 0; i < int(common.FoodNum); i++ {
			//并未被吃
			mFoods.eatfood[i] = false
			x, y := common.RandPOINTFloat64()
			mFoods.foodlist = append(mFoods.foodlist, common.Food{
				Energy: common.FoodEnergy,
				Stat:   common.POINT{X: x - 10, Y: y - 10},
			})
		}
	}
	//食物未生成够 添加食物
	if len(mFoods.foodlist) < cap(mFoods.foodlist) {
		//当容量小于 最大容量
		for i := 0; i < int(common.FoodNum); i++ {
			//食物已经被吃了 生成新food
			if mFoods.eatfood[i] {
				x, y := common.RandPOINTFloat64()
				mFoods.foodlist[i] = common.Food{
					Energy: common.FoodEnergy,
					Stat:   common.POINT{X: x - 10, Y: y - 10},
				}
				mFoods.eatfood[i] = false
			}
		}
	}
}

func (this *Scene) GetFoodList() *FoodList {
	return mFoods
}

//碰撞检测
func (this *Scene) CollisionDetection() bool {
	//撞墙检测
	if this.WallCollision() {
		this.snake.SnakeDie()
		return true
	}

	//吃食物
	for i, v := range mFoods.foodlist {
		if this.EatJuge(i, v) {
			mFoods.eatfood[i] = true
			this.EatFood(v)
		}
	}

	return false
}

// WallCollision 以下为碰撞检测部分
//撞墙
func (this *Scene) WallCollision() bool {

	if this.snake.head.X-this.snake.radius <= 0 ||
		this.snake.head.X+this.snake.radius >= common.SceneWidth ||
		this.snake.head.Y-this.snake.radius <= 0 ||
		this.snake.head.Y+this.snake.radius >= common.SceneHeight {
		return true
	}
	return false
}

// SnakeCollisionJudge 撞人
func (this *Scene) SnakeCollisionJudge() bool {

	head := this.snake.head

	var minD, headR, bodyR float64
	headR = this.snake.radius
	//todo invincible

	//遍历所有其他snakebody 的蛇身
	for _, othersnake := range this.others {
		bodyR = othersnake.radius
		minD = headR + bodyR

		//检测头部是否碰撞
		headX := head.X - othersnake.head.X
		headY := head.Y - othersnake.head.Y
		headD := math.Sqrt(headX*headX + headY*headY)
		//头头相撞 小的死
		if headD <= minD || len(this.snake.body) <= len(othersnake.body) {
			return true
		}

		//检测头部与身体是否碰撞
		for _, body := range othersnake.body {
			temX := head.X - body.X
			temY := head.Y - body.Y
			d := math.Sqrt(temX*temX + temY*temY)
			if d <= minD {
				return true
			}
		}

	}

	return false
}

// EatJuge 撞食物
func (this *Scene) EatJuge(index int, food common.Food) bool {

	temX := this.snake.head.X - food.Stat.X
	temY := this.snake.head.Y - food.Stat.Y
	d := math.Sqrt(temX*temX + temY*temY)

	//食物没被吃 并且到达了被吃的范围
	if !mFoods.eatfood[index] && d <= float64(food.Energy)+this.snake.radius+common.EatFoodRadius {
		return true
	}

	glog.Info("Eat Food")
	return false
}

func (this *Scene) EatFood(food common.Food) {

	this.snake.score = this.snake.score + food.Energy

	if len(this.snake.body) >= 300 {
		return
	}

	for i := 0; i < int(food.Energy); i++ {
		this.snake.body = append(this.snake.body, this.snake.body[len(this.snake.body)-1])
	}
}

func (this *Scene) InitSnake() {

	headx, heady := common.RandPOINTFloat64()
	temhead := common.POINT{
		X: headx - 50,
		Y: heady - 50,
	}

	this.snake.id = this.snake.thisplayer.id
	this.snake.name = this.snake.thisplayer.name
	this.snake.direct = this.snake.thisplayer.angle

	this.snake.head = temhead //初始化新生成的头
	//身体 默认向右移动
	for i := 1; i <= 3; i++ {
		this.snake.body = append(this.snake.body, common.POINT{
			X: temhead.X - 2*float64(i)*common.SnakeRadius,
			Y: temhead.Y,
		})
	}

	this.snake.score = 0
	this.snake.radius = 10
	//this.snake.invincible = true //无敌

}

func (this *Scene) UpdateSnakePOINT(angle uint32) {
	space := common.SceneSpeed * float64(((time.Now().UnixNano()/1e6)-(this.preTime))/1000) //相差(毫秒 / 1000) 即每秒
	this.preTime = (time.Now().UnixNano() / 1e6)                                            //毫秒

	this.SnakeMove(angle, space)
}

//todo 传的是引用可以吗?
func (this *Scene) UpdateOthersSnake() {

	for _, p := range this.room.players {
		if p.scene.others == nil {
			p.scene.others = []SnakeBody{}
		}

		for _, other := range this.room.players {
			//如果不是本身的那一条连接
			if p.id != other.id {
				p.scene.others = append(p.scene.others, SnakeBody{
					id:         other.scene.snake.id,
					name:       other.scene.snake.name,
					thisplayer: other.scene.snake.thisplayer,
					direct:     0,
					head:       other.scene.snake.head,
					body:       other.scene.snake.body,
				})

			}
		}
	}

}

func (this *Scene) SnakeMove(angle uint32, space float64) {
	//todo 算出蛇头移动后的坐标
	newhead := common.POINT{
		X: 0,
		Y: 0,
	}

	//碰撞检测 SnakBodyMove
	this.SnakeBodyMove(newhead)
	this.CollisionDetection()
}

func (this *Scene) SnakeBodyMove(newhead common.POINT) {
	//todo 计算出单位时间内移动的距离算出
	//判断是否吃食物

}

//todo 转化为食物 断开连接
func (this *Scene) SnakeDie(snake SnakeBody) {

	snake.thisplayer.wstask.Close()
}

func (this *Scene) SceneMsg() []byte {

	//序列化场景信息 本条蛇 其他玩家 食物数组
	var retsceneMsg common.RetSceneMsg
	//先分配空间
	retsceneMsg.PlayerSnake = []common.RetSnakeBody{}
	retsceneMsg.OthersSnake = []common.RetSnakeBody{}
	retsceneMsg.FoodList = []common.Food{}

	//本条蛇
	retsceneMsg.PlayerSnake = append(retsceneMsg.PlayerSnake, common.RetSnakeBody{
		Id:   this.snake.id,
		Name: this.snake.name,
		Head: this.snake.head,
		Body: this.snake.body,
	})
	//其他蛇
	for _, other := range this.others {
		retsceneMsg.OthersSnake = append(retsceneMsg.OthersSnake, common.RetSnakeBody{
			Id:   other.id,
			Name: other.name,
			Head: other.head,
			Body: other.body,
		})
	}
	//食物
	for i := 0; i <= int(common.FoodNum); i++ {
		//发送所有没被吃的食物
		if !mFoods.eatfood[i] {
			retsceneMsg.FoodList = append(retsceneMsg.FoodList, common.Food{
				Energy: mFoods.foodlist[i].Energy,
				Stat:   mFoods.foodlist[i].Stat,
			})
		}
	}

	bytes, err := json.Marshal(retsceneMsg)
	if nil != err {
		glog.Error("[Scene] Scene Msg Error ", err)
		return nil
	}

	return bytes
}
