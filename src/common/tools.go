package common

import (
	"math"
	"math/rand"
)

//求两点之间距离
func TwoPointLen(a, b POINT) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

//求三角形外接圆
func CircleOfTriangle(a, b, c POINT) Circle {
	var (
		a1 = 2 * (b.X - a.X)
		b1 = 2 * (b.Y - a.Y)
		c1 = b.X*b.X + b.Y*b.Y - a.X*a.X - a.Y*a.Y
		a2 = 2 * (c.X - b.Y)
		b2 = 2 * (c.Y - b.Y)
		c2 = c.X*c.X + c.Y*c.Y - b.X*b.X - b.Y*b.Y
	)

	center := POINT{
		X: (c1*b2 - c2*b1) / (a1*b2 - a2*b1),
		Y: (a1*c2 - a2*c1) / (a1*b2 - a2*b1),
	}

	return Circle{
		center: center,
		radius: TwoPointLen(a, center),
	}
}

//求最小覆盖圆
func MinCircle(pArr []POINT) Circle {
	temO := Circle{
		center: pArr[0],
		radius: 0,
	}

	for i := 0; i < len(pArr); i++ {
		if TwoPointLen(pArr[i], temO.center) > temO.radius {
			temO.center = pArr[i]
			temO.radius = 0
		}

		for j := 0; j < i; j++ {
			if TwoPointLen(pArr[i], temO.center) > temO.radius {
				temO.center = POINT{
					X: (pArr[i].X + pArr[j].X) / 2,
					Y: (pArr[i].Y + pArr[j].Y) / 2,
				}
				temO.radius = TwoPointLen(pArr[i], pArr[j]) / 2

				for k := 0; k < j; k++ {
					if TwoPointLen(pArr[k], temO.center) > temO.radius {
						temO = CircleOfTriangle(pArr[i], pArr[j], pArr[k])
					}
				}
			}
		}
	}
	return temO
}

type Snake struct {
	Name   string  //蛇名字
	Head   POINT   //蛇头
	Body   []POINT //蛇身数组
	Alive  bool    //是否存活
	Radius float64 //蛇的半径

	Color     string //蛇身颜色
	HeadColor string //蛇头颜色
}

func RandColor() uint32 {
	//todo: 随机生成颜色
	return 1 + uint32(rand.Intn(255))
}
