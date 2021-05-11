package main

import (
	"common"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"math"
	"time"
)

type SnakeBody struct {
	id         uint32
	name       string
	thisplayer *PlayerTask
	direct     float64
	head       common.POINT
	body       []common.POINT
	score      int32
	radius     float64
	isdead     bool
	//invincible bool todo 无敌
}

func (this *SnakeBody) SnakeDie() {
	fmt.Println("id: ", this.id, "name: ", this.id, "die")
	glog.Info("[snake die]")
	this.isdead = true //死了
	this.thisplayer.room.curnum--
	delete(this.thisplayer.room.players, this.thisplayer.id)
	this.thisplayer.OnClose()
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
	//others []SnakeBody //其他人的信息

	speed   float64 //蛇的移动速度
	preTime int64   //上一帧时间
}

func (this *Room) AddFoods() {
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
		fmt.Println("生成食物数组")
		glog.Info("生成食物数组")
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
		fmt.Println("填满食物数组")
		glog.Info("填满食物数组")
	}
}

func (this *Room) GetFoodList() *FoodList {
	return mFoods
}

//碰撞检测
func (this *Scene) CollisionDetection() bool {
	//撞墙检测
	if this.WallCollision() {
		fmt.Println("[Snake Die] 玩家: ", this.snake.thisplayer.id, " 撞墙了")
		glog.Info("[Snake Die] 玩家: ", this.snake.thisplayer.id, " 撞墙了")
		//this.snake.SnakeDie()
		return true
	}

	//吃食物
	for i, v := range mFoods.foodlist {
		if this.EatJuge(i, v) {
			this.EatFood(v)
			mFoods.eatfood[i] = true
		}
	}

	//撞人
	if this.SnakeCollisionJudge() {
		fmt.Println("[Snake Die] 玩家:", this.snake.thisplayer.id, "撞到了玩家")
		glog.Info("[Snake Die] 玩家:", this.snake.thisplayer.id, "撞到了玩家")
		//this.snake.SnakeDie()
		return true
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

	var minD float64

	//todo invincible

	//遍历所有其他 player 的 蛇身 蛇头
	for _, otherplayer := range this.room.players {
		//不是自身的那条链接
		if otherplayer != this.snake.thisplayer {

			fmt.Println(otherplayer)
			if otherplayer.scene == nil {
				continue
			}
			othersnake := &otherplayer.scene.snake

			minD = this.snake.radius + othersnake.radius

			//检测头部是否碰撞
			headX := this.snake.head.X - othersnake.head.X
			headY := this.snake.head.Y - othersnake.head.Y
			headD := math.Sqrt(headX*headX + headY*headY)

			//头头相撞 小的死
			if headD <= minD && len(this.snake.body) < len(othersnake.body) {
				fmt.Println("[Snake] 玩家ID: ", this.snake.id, "[蛇头位置: ]", " {", this.snake.head.X, " , ", this.snake.head.Y, "} ",
					"撞到了另一玩家ID: ", othersnake.id, "碰撞的蛇头位置为: ", " {", othersnake.head.X, " , ", othersnake.head.Y, "} ")
				glog.Info("[Snake] 玩家ID: ", this.snake.id, "[蛇头位置: ]", " {", this.snake.head.X, " , ", this.snake.head.Y, "} ",
					"撞到了另一玩家ID: ", othersnake.id, "碰撞的蛇头位置为: ", " {", othersnake.head.X, " , ", othersnake.head.Y, "} ")
				return true
			}

			//检测头部与身体是否碰撞
			for _, body := range othersnake.body {
				minD = this.snake.radius + othersnake.radius
				temX := this.snake.head.X - body.X
				temY := this.snake.head.Y - body.Y
				d := math.Sqrt(temX*temX + temY*temY)
				if d <= minD {
					fmt.Println("[Snake] 玩家ID: ", this.snake.id, "[蛇头位置: ]", " {", this.snake.head.X, " , ", this.snake.head.Y, "} ",
						"撞到了另一玩家ID: ", othersnake.id, "碰撞的蛇身位置为: ", " {", body.X, " , ", body.Y, "} ")
					glog.Info("[Snake] 玩家ID: ", this.snake.id, "[蛇头位置: ]", " {", this.snake.head.X, " , ", this.snake.head.Y, "} ",
						"撞到了另一玩家ID: ", othersnake.id, "碰撞的蛇身位置为: ", " {", body.X, " , ", body.Y, "} ")
					return true
				}
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
		glog.Info("Should Eat Food")
		return true
	}

	return false
}

func (this *Scene) EatFood(food common.Food) {
	fmt.Println("[Snake Eat] 玩家: ", this.snake.thisplayer.id, " 吃食物")
	glog.Info("[Snake Eat] 玩家: ", this.snake.thisplayer.id, " 吃食物")

	this.snake.score = this.snake.score + food.Energy

	if len(this.snake.body) >= 300 {
		return
	}

	for i := 0; i < int(food.Energy); i++ {
		this.snake.body = append(this.snake.body, this.snake.body[len(this.snake.body)-1])
	}
}

func (this *Scene) InitSnake() {

	headx, heady := common.SafeRandHeadFloat64(100, common.SceneWidth-100, common.SceneHeight-100)
	temhead := common.POINT{
		X: headx,
		Y: heady,
	}

	this.snake.head = temhead //初始化新生成的头
	this.snake.id = this.snake.thisplayer.id
	this.snake.name = this.snake.thisplayer.name
	this.snake.direct = this.snake.thisplayer.angle
	this.snake.radius = common.SnakeRadius
	this.snake.score = 0
	this.snake.isdead = false //没死
	//this.snake.invincible = true //无敌
	fmt.Println("玩家ID: ", this.snake.id, " [新初始化蛇头位置:]", "{", temhead.X, ",", temhead.Y, "}", "玩家姓名:", this.snake.name)
	glog.Info("玩家ID: ", this.snake.id, " [新初始化蛇头位置:]", "{", temhead.X, ",", temhead.Y, "}", "玩家姓名:", this.snake.name)

	//身体 默认向右移动 todo移动坐标有误
	for i := 1; i <= 3; i++ {
		this.snake.body = append(this.snake.body, common.POINT{
			X: temhead.X - float64(i),
			Y: temhead.Y,
		})
	}

	fmt.Println("蛇身位置", this.snake.body)
}

func (this *Scene) UpdateSnakePOINT(angle float64) {
	//space := common.SceneSpeed * float64(((time.Now().UnixNano()/1e6)-(this.preTime))/1000) //相差(毫秒 / 1000) 即每秒
	//this.preTime = (time.Now().UnixNano() / 1e6)                                            //毫秒
	//
	//this.SnakeMove(angle, space)

	frame := float64((time.Now().UnixNano() / 1e6) - (this.preTime)) //相差多少毫秒
	if frame > common.FrameTime {
		frame = common.FrameTime
	}

	//space := common.SceneSpeed * (frame / 1000) //速度像素/s 相差(毫秒 / 1000)即每秒

	this.preTime = time.Now().UnixNano() / 1e6 //毫秒

	this.SnakeHeadMove(angle, frame)
}

////todo 传的是引用可以吗?
//func (this *Scene) UpdateOthersSnake() {
//
//	for _, p := range this.room.players {
//		//if p.scene.others == nil {
//		//	p.scene.others = []SnakeBody{}
//		//}
//
//		for _, other := range this.room.players {
//			//如果不是本身的那一条连接
//			if p.id != other.id {
//				p.scene.others = append(p.scene.others, SnakeBody{
//					id:         other.scene.snake.id,
//					name:       other.scene.snake.name,
//					thisplayer: other.scene.snake.thisplayer,
//					direct:     0,
//					head:       other.scene.snake.head,
//					body:       other.scene.snake.body,
//				})
//
//			}
//		}
//	}
//
//}

func (this *Scene) UpdateSpeed(Speed float64) {
	this.speed = Speed
}

func (this *Scene) SnakeHeadMove(angle float64, space float64) {
	//蛇身越长 蛇的半径越大
	temRadius := math.Floor(float64(12 + len(this.snake.body)/100))
	this.snake.radius = temRadius
	//蛇身越大 速度也需要改变
	//this.speed = 105 + 400 / this.snake.radius

	//根据转向角度计算 蛇头朝向
	if math.Abs(angle-this.snake.direct) < 180 {
		if angle-this.snake.direct > 0 {
			if common.SnakeTurnSpeed*space > angle-this.snake.direct {
				this.snake.direct = angle
			} else {
				this.snake.direct += common.SnakeTurnSpeed * space
			}
		} else if angle-this.snake.direct < 0 {
			if common.SnakeTurnSpeed*space > this.snake.direct-angle {
				this.snake.direct = angle
			} else {
				this.snake.direct -= common.SnakeTurnSpeed * space
			}
		}
	}

	if math.Abs(angle-this.snake.direct) > 180 {
		if angle-this.snake.direct > 0 {
			this.snake.direct -= common.SnakeTurnSpeed * space
			if this.snake.direct < 0 {
				this.snake.direct += 360
				if this.snake.direct < angle {
					this.snake.direct = angle
				}
			}
		} else if angle-this.snake.direct < 0 {
			this.snake.direct += common.SnakeTurnSpeed * space

			if this.snake.direct > 360 {
				this.snake.direct -= 360
				if this.snake.direct > angle {
					this.snake.direct = angle
				}
			}
		}
	}

	moveX := this.speed * space / 1000 * math.Cos(math.Pi*this.snake.direct/180)
	moveY := this.speed * space / 1000 * math.Sin(math.Pi*this.snake.direct/180)
	//moveDistance := this.speed * space / 1000
	newhead := common.POINT{
		X: this.snake.head.X + moveX,
		Y: this.snake.head.Y - moveY,
	}
	//this.snake.movedistance += moveDistance
	//碰撞检测 SnakBodyMove
	fmt.Println("[Snake Move] 玩家ID: ", this.snake.id, " [新生成蛇头位置:]", "{", newhead.X, ",", newhead.Y, "}", "玩家姓名:", this.snake.name)
	glog.Info("[Snake Move] 玩家ID: ", this.snake.id, " [新生成蛇头位置:]", "{", newhead.X, ",", newhead.Y, "}", "玩家姓名:", this.snake.name)
	this.SnakeBodyMove(newhead)

	fmt.Println(this.snake.body)

	dead := this.CollisionDetection()

	if dead {
		this.snake.SnakeDie()
	}

	this.snake.head.X += moveX
	this.snake.head.Y -= moveY
}

func (this *Scene) SnakeBodyMove(newhead common.POINT) {
	//todo 计算出单位时间内移动的距离算出
	for i := len(this.snake.body) - 1; i > 0; i-- {
		this.snake.body[i] = this.snake.body[i-1]
	}
	this.snake.body[0] = this.snake.head
	this.snake.head = newhead

}

//todo 转化为食物 断开连接
func (this *Scene) SnakeDie(snake SnakeBody) {

	snake.thisplayer.wstask.Close()
}

func (this *Scene) SceneMsg() []byte {

	if this.snake.isdead {
		return nil
	}

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
	for _, other := range this.room.players {
		if other != this.snake.thisplayer {
			fmt.Println("玩家: ", this.snake.id, "[序列化] 不属于本条蛇的信息: ", other)
			retsceneMsg.OthersSnake = append(retsceneMsg.OthersSnake, common.RetSnakeBody{
				Id:   other.scene.snake.id,
				Name: other.scene.snake.name,
				Head: other.scene.snake.head,
				Body: other.scene.snake.body,
			})
		}

	}

	//食物
	for i := 0; i < int(common.FoodNum); i++ {
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
		glog.Error("[Scene] SceneMsg 序列化出错 ", err)
		return nil
	}

	return bytes
}
