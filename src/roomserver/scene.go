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

	otherhead []common.POINT //其他人的头
	otherbody []common.POINT //其他人的身体 做判断
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
