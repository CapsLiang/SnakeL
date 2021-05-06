package main

import (
	"common"
	"encoding/json"
	"github.com/golang/glog"
	"time"
)

type SnakeBody struct {
	id         uint32
	name       string
	thisplayer *PlayerTask
	direct     uint32
	head       common.POINT
	body       []common.POINT
	//invincible bool todo 无敌
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

func (this *Scene) EatFood() {
	//todo this.snake.Head
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
	//this.snake.invincible = true //无敌

}

func (this *Scene) UpdateSnakePOINT(angle uint32) {
	space := common.SceneSpeed * float64(((time.Now().UnixNano()/1e6)-(this.preTime))/1000) //相差(毫秒 / 1000) 即每秒
	this.preTime = (time.Now().UnixNano() / 1e6)                                            //毫秒

	this.SnakeMove(angle, space)
}

//todo 传的是引用吗?
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
	//SnakBodyMove
	this.SnakeBodyMove(newhead)
}

func (this *Scene) SnakeBodyMove(newhead common.POINT) {
	//todo 计算出单位时间内移动的距离算出
}

func (this *Scene) SceneMsg() []byte {

	//todo 序列化场景信息
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
