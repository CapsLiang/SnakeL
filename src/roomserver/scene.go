package main

import (
	"common"
	"encoding/json"
	"github.com/golang/glog"
)

type Scene struct {
	room *Room

	sceneWidth     float64
	sceneHeight    float64
	sceneGridSize  float64
	sceneGridColor string

	head       common.POINT   //蛇头的位置
	body       []common.POINT //蛇身
	headnext   common.POINT   //蛇头下次移动的位置
	invincible bool           //无敌
	//pool *BallPool //后续做不同类型的食物

	foodlist []common.POINT //食物列表
	eatfood  map[int]bool   //判断是否被吃 每次更新

	otherhead []common.POINT //其他人的头
	otherbody []common.POINT //其他人的身体 做判断
}

func (this *Scene) FoodList_GetMe() []common.POINT {
	if nil == this.foodlist {
		this.foodlist = make([]common.POINT, 0, common.FoodNum)
		//初始化食物数组
		for i := 0; i <= int(common.FoodNum); i++ {
			//并未被吃
			this.eatfood[i] = false
			x, y := common.RandPOINTFloat64()
			this.foodlist = append(this.foodlist, common.POINT{
				X: x - 10,
				Y: y - 10,
			})
		}

	}
	return this.foodlist
}

//添加食物
func (this *Scene) AddFoods() {
	if len(this.foodlist) < cap(this.foodlist) {
		//当容量小于 最大容量
		for i := 0; i <= int(common.FoodNum); i++ {
			//食物已经被吃了 生成新food
			if !this.eatfood[i] {
				this.eatfood[i] = true
				x, y := common.RandPOINTFloat64()
				this.foodlist[i] = common.POINT{
					X: x - 10,
					Y: y - 10,
				}
			}
		}
	}
}

func (this *Scene) EatFood() {

}

func (this *Scene) InitSnake() {
	headx, heady := common.RandPOINTFloat64()
	temhead := common.POINT{
		X: headx - 20,
		Y: heady - 20,
	}

	this.head = temhead

	for i := 1; i <= 3; i++ {
		this.body = append(this.body, common.POINT{
			X: temhead.X - 2*float64(i)*common.SnakeRadius,
			Y: temhead.X - 2*float64(i)*common.SnakeRadius,
		})
	}

}

func (this *Scene) SnakeMove(angle, space uint32) {
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
	retsceneMsg = common.RetSceneMsg{}

	bytes, err := json.Marshal(retsceneMsg)
	if nil != err {
		glog.Error("[Scene] Scene Msg Error ", err)
		return nil
	}

	return bytes
}
