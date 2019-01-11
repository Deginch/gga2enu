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
			//logE("==========")
			//logE("gga=%v",gga)
			//logE("gga=%v",oneEnu)
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
	if(1==2){
		gps2enuMain()
		testGga2Enu()
		return
	}
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

func testGga2Enu(){
	//var ggaStr="$GPRMC,061807.00,A,4001.4906698,N,11617.2357373,E,0.05,0.00,060119,0.0,E,D*3F\n$GPGGA,061807.00,4001.4906698,N,11617.2357373,E,4,16,1.0,65.239,M,-9.546,M,1.0,*6D\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,34,051,46,14,58,317,46,20,19,171,40,1*69\n$GPGSV,2,2,07,25,58,091,47,31,41,255,47,32,75,354,49,,,,,1*58\n$GPRMC,061808.00,A,4001.4906244,N,11617.2356983,E,0.04,0.00,060119,0.0,E,D*30\n$GPGGA,061808.00,4001.4906244,N,11617.2356983,E,4,16,1.0,65.607,M,-9.546,M,1.0,*6A\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,34,051,45,14,58,317,46,20,19,171,40,1*6A\n$GPGSV,2,2,07,25,58,091,46,31,41,255,47,32,75,354,49,,,,,1*59\n$GPRMC,061809.00,A,4001.4906755,N,11617.2357374,E,0.05,0.00,060119,0.0,E,D*36\n$GPGGA,061809.00,4001.4906755,N,11617.2357374,E,4,16,1.0,65.594,M,-9.546,M,1.0,*64\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,34,051,46,14,58,317,46,20,19,171,40,1*69\n$GPGSV,2,2,07,25,58,091,46,31,41,255,47,32,75,354,49,,,,,1*59\n$GPRMC,061810.00,A,4001.4907937,N,11617.2356564,E,0.03,0.00,060119,0.0,E,D*35\n$GPGGA,061810.00,4001.4907937,N,11617.2356564,E,4,16,1.0,65.298,M,-9.546,M,1.0,*6A\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,33,051,46,14,58,317,46,20,19,171,40,1*6E\n$GPGSV,2,2,07,25,58,091,46,31,41,255,47,32,75,354,49,,,,,1*59\n$GPRMC,061811.00,A,4001.4908099,N,11617.2356831,E,0.09,0.00,060119,0.0,E,D*31\n$GPGGA,061811.00,4001.4908099,N,11617.2356831,E,4,16,1.0,65.116,M,-9.546,M,1.0,*61\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,33,051,46,14,58,317,46,20,19,171,40,1*6E\n$GPGSV,2,2,07,25,58,091,46,31,41,255,47,32,75,354,49,,,,,1*59\n$GPRMC,061812.00,A,4001.4907875,N,11617.2356365,E,0.04,0.00,060119,0.0,E,D*30\n$GPGGA,061812.00,4001.4907875,N,11617.2356365,E,4,16,1.0,65.206,M,-9.546,M,1.0,*6F\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,33,051,46,14,58,317,46,20,19,171,40,1*6E\n$GPGSV,2,2,07,25,58,091,46,31,41,255,47,32,75,354,49,,,,,1*59\n$GPRMC,061813.00,A,4001.4907372,N,11617.2356211,E,0.04,0.00,060119,0.0,E,D*3F\n$GPGGA,061813.00,4001.4907372,N,11617.2356211,E,4,16,1.0,65.331,M,-9.546,M,1.0,*65\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,33,051,46,14,58,317,46,20,19,171,40,1*6E\n$GPGSV,2,2,07,25,58,091,46,31,41,255,46,32,75,354,49,,,,,1*58\n"
	var ggaStr2="$GPGGA,152534.00,4001.4901221,N,11617.2360899,E,4,19,1.0,70.712,M,-9.546,M,1.1,*6D"
	rlat := 40.024813632
	rlon := 116.287156429
	rh := 61.384523
	ggaArray,_:=parseGGA(ggaStr2)
	for _,gga:=range ggaArray{
		enu,_:=gga.toENU(rlat,rlon,rh)
		fmt.Println("=========")
		fmt.Println(gga)
		fmt.Println(enu)
	}

}