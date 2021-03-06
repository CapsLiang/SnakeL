package common

/***************/
/*****Http******/
/***************/

// RetGetIDMsg 服务器返回id
type RetGetIDMsg struct {
	Id uint32 `json:"id"`
}

// ReqGetIDMsg 客户请求 json 设备信息和ip地址 存在redis中
type ReqGetIDMsg struct {
	DeviceId string `json:"deviceId"`
	Ip       string `json:"ip"`
}

// RetOverMsg 发送结束消息
type RetOverMsg struct {
	End bool `json:"end"`
}

type RetTimeMsg struct {
	Time uint64 `json:"time"`
}

// RetSnakeBody omitempty作用是在json数据结构转换时，当该字段的值为该字段类型的零值时，忽略该字段
type RetSnakeBody struct {
	//color
	Id     uint32  `json:"id"`
	Name   string  `json:"name"`
	Direct float64 `json:"direct"`
	Head   POINT   `json:"head"`
	Body   []POINT `json:"body"`
	Score  int32   `json:"score"`
	Radius float64 `json:"radius"`
}

// RetSceneMsg 服务器返回游戏场景
type RetSceneMsg struct {
	PlayerSnake []RetSnakeBody `json:"player_snake"`
	OthersSnake []RetSnakeBody `json:"others_snake"`
	FoodList    []Food         `json:"food_list"`
}

// MsgType 消息信息
type MsgType uint16

const (
	MsgType_Start   MsgType = 0
	MsgType_Move    MsgType = 1
	MsgType_SpeedUp MsgType = 2
	MsgType_Finsh   MsgType = 3
	MsgType_Heart   MsgType = 4
)

//连接信息
const (
	// Task_TimeOut websocketTask 超时时间
	Task_TimeOut = 20
)

//场景信息 速度 宽度 高度
const (
	// FrameTime 20ms
	FrameTime float64 = 20
	// SceneSpeed 像素每秒
	SceneSpeed float64 = 200
	// SceneSpeedUp 加速系数
	SceneSpeedUp float64 = 1.6
	// SceneWidth 场景的宽度
	SceneWidth float64 = 800
	// SceneHeight 场景高度
	SceneHeight float64 = 800
)

//食物信息
const (
	FoodNum    uint32 = 20
	FoodEnergy int32  = 5
	// EatFoodRadius
	EatFoodRadius float64 = 2
)

//蛇的信息
const (
	// SnakeRadius 蛇的半径 10px
	SnakeRadius float64 = 10
	// SnakeTurnSpeed 蛇的转向速度 弧度/ms
	SnakeTurnSpeed = 0.5
)

//以下为约定结构
type POINT struct {
	//Id uint32
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Food struct {
	Energy int32 `json:"energy"`
	Stat   POINT `json:"stat"`
}

type Circle struct {
	center POINT
	radius float64
}
