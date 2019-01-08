package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const GPGGA = "$GPGGA"


type GGA struct {
	/*
	utf时间 hhmmss.sss
	 */
	utc string
	/*
	维度
	 */
	lat float64
	/*
	n/s指示，N北纬，S南纬
	 */
	ns string
	/*
	经度
	 */
	lon float64
	/*
	EW指示，E东经，W西经
	 */
	ew string
	/*
	质量，0未定位，1gps单点，2差分定位，3pps解，4rtk固定解，5rtk浮点解，6估计值，7手工输入，8模拟模式
	 */
	quality int
	/*
	卫星数量
	 */
	satelliteNum int
	hdop         float64
	xHeight         float64
	/*
	高度单位
	 */
	heightUnit string
	/*
	高度
	 */
	height float64

	/*
	原始数据
	 */
	raw string
}

var qualityMap = map[int]int{
	4: 1,
	5: 2,
	1: 5,
	2: 1,
}
/*
将gga转换成enu，需要传入基准经纬高
 */
func (gga GGA) toENU(lat float64, lon float64, height float64) (ENU, error) {
	var enu ENU;
	//region 默认数据
	enu.date=fmt.Sprintf("%s %s:%s:%s",time.Now().Format("2006-01-02"),gga.utc[:2],gga.utc[2:4],gga.utc[4:])
	enu.quality = gga.quality
	enu.satelliteNum = gga.satelliteNum
	enu.jetLag = 0
	//endregion
	//region 开始处理质量
	enuQuality := qualityMap[gga.quality]
	if enuQuality == 0 {
		enuQuality = gga.quality
	}
	enu.quality = enuQuality;
	//endregion
	//region 处理置信率
	switch enu.quality {
	case 1:
		enu.confidenceRate = 999.9
	case 2:
		enu.confidenceRate = 1.0
	default:
		enu.confidenceRate = 0
	}
	//endregion
	//region 处理经纬高
	var err error;
	enu.east, enu.north, enu.up = calGps2Enu(lat,lon,height,gga.lat,gga.lon,gga.xHeight+gga.height)
	if err != nil {
		return enu, err
	}
	//endregion
	return enu, nil
}

type ENU struct {
	date  string
	east  float64
	north float64
	up    float64
	/*
	质量，0未定位，1gps单点，2差分定位，3pps解，4rtk固定解，5rtk浮点解，6估计值，7手工输入，8模拟模式
	 */
	quality int
	/*
	卫星数量
	 */
	satelliteNum int
	other1       float64
	other2       float64
	other3       float64
	other4       float64
	other5       float64
	other6       float64
	/*
	时差
	 */
	jetLag float64
	/*
	置信率
	 */
	confidenceRate float64
}

func (enu ENU) toString() string {
	return fmt.Sprintf("%s %f %f %f %d %d %f %f %f %f %f %f %f %f\n",
		enu.date,
		enu.east,
		enu.north,
		enu.up,
		enu.quality,
		enu.satelliteNum,
		enu.other1,
		enu.other2,
		enu.other3,
		enu.other4,
		enu.other5,
		enu.other6,
		enu.jetLag,
		enu.confidenceRate,
	)
}

func parseGGA(data string) ([]GGA, string) {
	var cache = ""
	dataArray := strings.Split(data, "\n");
	ggaArray := make([]GGA, 0)
	for _, line := range dataArray {
		eleArray := strings.Split(line, ",")
		if len(eleArray) != 15 {
			cache = line
			continue
		}
		if eleArray[0] != GPGGA {
			continue
		}
		var gga GGA;
		gga.utc= eleArray[1]
		gga.lat, _ = strconv.ParseFloat(eleArray[2], 64)
		gga.ns = eleArray[3]
		gga.lon, _ = strconv.ParseFloat(eleArray[4], 64)
		gga.ew = eleArray[5]
		gga.lat=angle2decimal(gga.lat/100)
		gga.lon=angle2decimal(gga.lon/100)
		gga.quality, _ = strconv.Atoi(eleArray[6])
		gga.satelliteNum, _ = strconv.Atoi(eleArray[7])
		gga.hdop, _ = strconv.ParseFloat(eleArray[8], 64)
		gga.xHeight,_=strconv.ParseFloat(eleArray[9],64)
		gga.heightUnit = eleArray[10]
		gga.height, _ = strconv.ParseFloat(eleArray[11], 64)
		gga.raw = line
		ggaArray = append(ggaArray, gga)
	}
	return ggaArray, cache
}

func angle2decimal(angle float64)float64{
	left:=int(angle)
	right:=angle-float64(left)
	return float64(left)+right*100/60
}

var testGGA = "$GPGGA,152534.00,4001.4901221,N,11617.2360899,E,4,19,1.0,70.712,M,-9.546,M,1.1,*6D"

var testENU = "2019/01/06 06:18:43.000         9.5356         2.4139        -0.2141   1  16   0.0049   0.0060   0.0134  -0.0025  -0.0072  -0.0006   1.00  999.9"

func __main() {
	ggaArray, remain := parseGGA(testGGA)
	fmt.Println(ggaArray, remain)
}
