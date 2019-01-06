package main

import (
	"fmt"
	"github.com/im7mortal/UTM"
	"github.com/pkg/errors"
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
	enu.east, enu.north, err = latLon2EastNorth(gga.lat, gga.lon, gga.ns == "N", lat, lon, true)
	if err != nil {
		return enu, err
	}
	enu.height = gga.height - height
	//endregion
	return enu, nil
}

type ENU struct {
	date   string
	east   float64
	north  float64
	height float64
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
		enu.height,
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
		gga.quality, _ = strconv.Atoi(eleArray[6])
		gga.satelliteNum, _ = strconv.Atoi(eleArray[7])
		gga.hdop, _ = strconv.ParseFloat(eleArray[8], 64)
		gga.heightUnit = eleArray[9]
		gga.height, _ = strconv.ParseFloat(eleArray[10], 64)
		gga.raw = line
		ggaArray = append(ggaArray, gga)
	}
	return ggaArray, cache
}

func latLon2EastNorth(llat float64, llon float64, lNorthern bool, rlat float64, rlon float64, rNorthern bool) (float64, float64, error) {
	le, ln, _, _, lerr := UTM.FromLatLon(llat, llon, lNorthern)
	re, rn, _, _, rerr := UTM.FromLatLon(rlat, rlon, rNorthern)
	if lerr != nil || rerr != nil {
		return 0, 0, errors.New(fmt.Sprintf("计算失败 latLon2EastNorth error,lerr=%s,rerr=%s",lerr,rerr))
	}
	return le - re, ln - rn, nil
}

var testGGA = "$GPRMC,061807.00,A,4001.4906698,N,11617.2357373,E,0.05,0.00,060119,0.0,E,D*3F\n$GPGGA,061807.00,4001.4906698,N,11617.2357373,E,4,16,1.0,65.239,M,-9.546,M,1.0,*6D\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,34,051,46,14,58,317,46,20,19,171,40,1*69\n$GPGSV,2,2,07,25,58,091,47,31,41,255,47,32,75,354,49,,,,,1*58\n$GPRMC,061808.00,A,4001.4906244,N,11617.2356983,E,0.04,0.00,060119,0.0,E,D*30\n$GPGGA,061808.00,4001.4906244,N,11617.2356983,E,4,16,1.0,65.607,M,-9.546,M,1.0,*6A\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,34,051,45,14,58,317,46,20,19,171,40,1*6A\n$GPGSV,2,2,07,25,58,091,46,31,41,255,47,32,75,354,49,,,,,1*59\n$GPRMC,061809.00,A,4001.4906755,N,11617.2357374,E,0.05,0.00,060119,0.0,E,D*36\n$GPGGA,061809.00,4001.4906755,N,11617.2357374,E,4,16,1.0,65.594,M,-9.546,M,1.0,*64\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,34,051,46,14,58,317,46,20,19,171,40,1*69\n$GPGSV,2,2,07,25,58,091,46,31,41,255,47,32,75,354,49,,,,,1*59\n$GPRMC,061810.00,A,4001.4907937,N,11617.2356564,E,0.03,0.00,060119,0.0,E,D*35\n$GPGGA,061810.00,4001.4907937,N,11617.2356564,E,4,16,1.0,65.298,M,-9.546,M,1.0,*6A\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,33,051,46,14,58,317,46,20,19,171,40,1*6E\n$GPGSV,2,2,07,25,58,091,46,31,41,255,47,32,75,354,49,,,,,1*59\n$GPRMC,061811.00,A,4001.4908099,N,11617.2356831,E,0.09,0.00,060119,0.0,E,D*31\n$GPGGA,061811.00,4001.4908099,N,11617.2356831,E,4,16,1.0,65.116,M,-9.546,M,1.0,*61\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,33,051,46,14,58,317,46,20,19,171,40,1*6E\n$GPGSV,2,2,07,25,58,091,46,31,41,255,47,32,75,354,49,,,,,1*59\n$GPRMC,061812.00,A,4001.4907875,N,11617.2356365,E,0.04,0.00,060119,0.0,E,D*30\n$GPGGA,061812.00,4001.4907875,N,11617.2356365,E,4,16,1.0,65.206,M,-9.546,M,1.0,*6F\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,33,051,46,14,58,317,46,20,19,171,40,1*6E\n$GPGSV,2,2,07,25,58,091,46,31,41,255,47,32,75,354,49,,,,,1*59\n$GPRMC,061813.00,A,4001.4907372,N,11617.2356211,E,0.04,0.00,060119,0.0,E,D*3F\n$GPGGA,061813.00,4001.4907372,N,11617.2356211,E,4,16,1.0,65.331,M,-9.546,M,1.0,*65\n$GPGSA,A,3,10,12,14,20,25,31,32,,,,,,2.7,1.3,2.4,1*2F\n$GPGSV,2,1,07,10,46,184,47,12,33,051,46,14,58,317,46,20,19,171,40,1*6E\n$GPGSV,2,2,07,25,58,091,46,31,41,255,46,32,75,354,49,,,,,1*58\n"

var testENU = "2019/01/06 06:18:43.000         9.5356         2.4139        -0.2141   1  16   0.0049   0.0060   0.0134  -0.0025  -0.0072  -0.0006   1.00  999.9"

func testMain() {
	ggaArray, remain := parseGGA(testGGA)
	fmt.Println(ggaArray, remain)
	fmt.Println("========")
	fmt.Println(latLon2EastNorth(40+1.4906698/60.0,116+17.2357373/60.0,true,40.024813632,116.287156429,true))
}
