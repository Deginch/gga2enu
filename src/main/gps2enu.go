package main

import (
	"fmt"
	"math"
)

var a = 6378137.0
var f = 1.0 / 298.257223563

var e2 = f * (2 - f)

func angle2R(angle float64) float64 {
	return angle / 180 * math.Pi
}

func calV(latR float64) float64 {
	sinLatR := math.Sin(latR)
	return a / (math.Sqrt(1 - e2*sinLatR*sinLatR))
}

func calRr(latR float64, lonR float64, h float64) []float64 {
	v := calV(latR)
	return []float64{
		(v + h) * math.Cos(latR) * math.Cos(lonR),
		(v + h) * math.Cos(latR) * math.Sin(lonR),
		v * (1 - e2) * math.Sin(latR),
	}
}

func calEr(latR float64, lonR float64) [][]float64 {
	return [][]float64{
		{
			- math.Sin(lonR), math.Cos(lonR), 0.0,
		},
		{
			-math.Sin(latR) * math.Cos(lonR), -math.Sin(latR) * math.Sin(lonR), math.Cos(latR),
		},
		{
			math.Cos(latR) * math.Cos(lonR), math.Cos(latR) * math.Sin(lonR), math.Sin(latR),
		},
	}
}
func vecMulVec(l []float64, r []float64) float64 {
	if len(l) != len(r) {
		return 0;
	}
	result := 0.0
	for i := range l {
		result += l[i] * r[i]
	}
	return result
}

func calEnu(Er [][]float64, rECEF []float64, rR []float64) (float64,float64,float64) {
	rTEMP := []float64{
		rECEF[0] - rR[0],
		rECEF[1] - rR[1],
		rECEF[2] - rR[2],
	}
	return vecMulVec(Er[0], rTEMP), vecMulVec(Er[1], rTEMP), vecMulVec(Er[2], rTEMP)
}

func calGps2Enu(rlat float64, rlon float64, rh float64, llat float64, llon float64, lh float64) (float64,float64,float64) {
	rECEF := calRr(angle2R(llat), angle2R(llon), lh)
	rR := calRr(angle2R(rlat), angle2R(rlon), rh)
	Er:=calEr(angle2R(rlat),angle2R(rlon))
	return calEnu(Er, rECEF, rR)
}

func gps2enuMain() {
	//llat:=40.02484449666667
	//llon:=116.28726228833334
	llat := 40.024835368333335
	llon := 116.287268165
	lh := 61.166
	rlat := 40.024813632
	rlon := 116.287156429
	rh := 61.384523
	//fmt.Println(calGps2Enu(llat,llon,lh,rlat,rlon,rh))
	fmt.Println(calGps2Enu(rlat, rlon, rh, llat, llon, lh))
}
