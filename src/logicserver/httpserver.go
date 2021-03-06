package main

import (
	"base/env"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	common "common"

	"github.com/golang/glog"
	"github.com/rs/cors"
)

func GetNameHandler(w http.ResponseWriter, r *http.Request) {
	randName := "NickName" + GetDateFormat()

	fmt.Fprintf(w, randName)
}

func GetIDHandler(w http.ResponseWriter, r *http.Request) {
	msg := common.ReqGetIDMsg{
		DeviceId: r.FormValue("DeviceId"),
		Ip:       r.FormValue("Ip"),
	}

	glog.Info("[login] 收到设备信息", msg)

	id, err := AccRedis_GetMe().GetIncID()
	if nil != err {
		glog.Error("[login] 获取自增ID失败 ", err)
		return
	}

	err = json.NewEncoder(w).Encode(&common.RetGetIDMsg{
		Id: id,
	})
	if nil != err {
		glog.Error("[login] 返回ID失败 ", err)
		return
	}

	err = AccRedis_GetMe().SetDeviceIdAndIp(id, &msg)
	if nil != err {
		glog.Error("[login] Set 用户信息失败 ", err)
		return
	}

}

func GetRoomHandler(w http.ResponseWriter, r *http.Request) {
	info, err := GetVailabelRoomInfo("token")
	if nil != err {
		glog.Error("[login] RPC获取房间信息失败 ", err)
		return
	}

	fmt.Fprintf(w, strconv.Itoa(int(info.Port)))
}

// 时间戳转年月日 时分秒
func GetDateFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func StartHttpServer() bool {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/getname", GetNameHandler)
	mux.HandleFunc("/getroom", GetRoomHandler)
	mux.HandleFunc("/getid", GetIDHandler)

	addr := env.Get("logic", "listen")
	handler := c.Handler(mux)
	http.ListenAndServe(addr, handler)

	glog.Info("[Start] Bind Port Success, Port=", addr)

	return true
}
