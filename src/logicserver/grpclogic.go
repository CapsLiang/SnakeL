package main

import (
	"base/env"
	"context"

	proto "proto"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

func GetVailabelRoomInfo(token string) (*proto.ConnectRoomInfo, error) {
	conn, err := grpc.Dial(env.Get("logic", "grpc"), grpc.WithInsecure())
	if nil != err {
		glog.Error("[login][gRPC] 连接失败 ", err)
		return nil, err
	}

	client := proto.NewStreamLogicServiceClient(conn)

	result, err := client.Route(context.Background(), &proto.LogicRequest{Token: token})
	if nil != err {
		glog.Error("[login][gRPC] 客户端调用错误 ", err)
		return nil, err
	}

	return result.MInfo, nil
}
