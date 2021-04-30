package main

//todo 对象池
//
//import "common"
//
//type BallFood struct {
//	balltype uint32
//	//todo color string
//	energy int32
//	center common.Circle
//	alive bool
//}
//
//type BallPool struct {
//	foodPool []*BallFood
//	foodIndex int
//}
//
//
//func NewBallPool() (pool *BallPool){
//	pool = &BallPool{}
//	for i:= 0; i < int(common.FoodPoolNum); i++ {
//		pool.foodIndex++
//		pool.foodPool = append(pool.foodPool, &BallFood{})
//	}
//	return
//}
//
//
//func (this *BallPool) GetFood() (food *BallFood) {
//	if this.foodIndex == 0 {
//		food = &BallFood{}
//	} else {
//		this.foodIndex--
//		food = this.foodPool[this.foodIndex]
//
//	}
//	return
//}
