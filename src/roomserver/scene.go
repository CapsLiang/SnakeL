package main

import (
	"common"
	"encoding/json"
	"github.com/golang/glog"
)

type Scene struct {
	sceneWidth     float64
	sceneHeight    float64
	sceneGridSize  float64
	sceneGridColor string

	head     common.POINT //蛇头的位置
	headnext common.POINT //蛇头下次移动的位置
	//pool *BallPool //后续做不同类型的食物

	foodlist []common.POINT //食物列表
	eatfood  map[int]bool   //判断是否被吃 每次更新

	otherhead []common.POINT //其他人的头
	otherbody []common.POINT //其他人的身体 做判断
}

func (this *Scene) FoodList_GetMe() []common.POINT {
	if nil == this.foodlist {
		this.foodlist = make([]common.POINT, 0, common.FoodNum)

		for i := 0; i <= int(common.FoodNum); i++ {
			this.eatfood[i] = true
			x, y := common.RandPOINTFloat64()
			this.foodlist = append(this.foodlist, common.POINT{
				X: x - 10,
				Y: y - 10,
			})
		}

	}
	return this.foodlist
}

func (this *Scene) AddFoods() {
	if len(this.foodlist) < cap(this.foodlist) {

	}
}

func InitScene() {

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
