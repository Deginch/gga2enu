package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func gga2enuHandle(lat float64,lon float64,height float64) handleData {
	var preCache=""
	var ggaArray []GGA
	return func(from []byte) []byte {
		ggaArray,preCache=parseGGA(preCache+string(from[:]))
		var enuByte =make([]byte,0)
		for _,gga:=range ggaArray{
			oneEnu,err:=gga.toENU(lat,lon,height)
			if err !=nil{
				logE("解析参数失败，当前gga为%v，错误原因为%v",gga,err)
				continue
			}
			enuByte=append(enuByte,[]byte(oneEnu.toString())...)
		}
		return enuByte
	}
}
func help() {
	fmt.Println("参数格式如下:")
	fmt.Println("ggaIp:ggaPort,enuIp:enuPort,经度,纬度,高度;ggaIp:ggaPort,enuIp:enuPort,经度,纬度,高度;ggaIp:ggaPort,enuIp:enuPort,经度,纬度,高度;")
	fmt.Println("例如:")
	fmt.Println("127.0.0.1:10000,127.0.0.10001,20.13,48.0,123.0;127.0.0.1:10000,127.0.0.10001,20.13,48.0,123.0;127.0.0.1:10000,127.0.0.10001,20.13,48.0,123.0;")
}

func main() {
	args := os.Args
	if len(args) < 2 || args == nil {
		help()
		return
	}
	all:=args[1]
	for _,one:=range strings.Split(all,";"){
		params:=strings.Split(one,",")
		if len(params)!=5{
			logE("参数错误，错误参数行为:%s",one)
			return
		}
		ipFrom:=params[0]
		ipTo:=params[1]
		lat,_:=strconv.ParseFloat(params[2],64)
		lon,_:=strconv.ParseFloat(params[3],64)
		height,_:=strconv.ParseFloat(params[4],64)
		go transfer(ipFrom,ipTo,3*time.Second,gga2enuHandle(lat,lon,height))
	}
	for{
		time.Sleep(10*time.Second)
	}
}
