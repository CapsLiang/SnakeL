package common

/***************/
/*****Http******/
/***************/

//服务器返回id
type RetGetIDMsg struct {
	Id uint32 `json:"id"`
}

// 客户请求 json 设备信息和ip地址 存在redis中
type ReqGetIDMsg struct {
	DeviceId string `json:"deviceId"`
	Ip       string `json:"ip"`
}

//发送结束消息
type RetOverMsg struct {
	End bool `json:"end"`
}

//omitempty作用是在json数据结构转换时，当该字段的值为该字段类型的零值时，忽略该字段
type RetSnakeBody struct {
	//color
	Id   uint32  `json:"id,omitempty"`
	Name string  `json:"name,omitempty"`
	Head POINT   `json:"head"`
	Body []POINT `json:"body,omitempty"`
}

//服务器返回游戏场景
type RetSceneMsg struct {
	PlayerSnake []RetSnakeBody `json:"player_snake"`
	OthersSnake []RetSnakeBody `json:"others_snake"`
	FoodList    []Food         `json:"food_list"`
}

//websocketTask 超时时间
const (
	Task_TimeOut = 20
)

//场景信息 速度 宽度 高度
const (
	SceneSpeed float64 = 200 //像素每秒
	//场景的高度与大小
	SceneWidth  float64 = 800
	SceneHeight float64 = 800
	//格子的颜色与大小
	SceneGridColor string  = "#f6f6f6"
	SceneGridSize  float64 = 20
)

//食物信息
const (
	FoodPoolNum uint32  = 2000
	FoodNum     uint32  = 20
	FoodEnergy  int32   = 5
	FoodRadius  float64 = 5
)

type POINT struct {
	//Id uint32
	X float64
	Y float64
}

type Food struct {
	Energy int32
	Stat   POINT
}

type Circle struct {
	center POINT
	radius float64
}

type MsgType uint8

const (
	//MsgType_Token  MsgType = 0
	//MsgType_Move   MsgType = 1
	//MsgType_Finsh  MsgType = 2
	//MsgType_Shoot  MsgType = 3
	//MsgType_Heart  MsgType = 4
	MsgType_Direct MsgType = 0 //传角度
)

//蛇
const (
	SnakeRadius float64 = 0.5
)

//建议用二进制传angle
//type UserCmd struct {
//	Id    uint32 `json:"id"`
//	Angle int32 `json:"angle"`
//}
