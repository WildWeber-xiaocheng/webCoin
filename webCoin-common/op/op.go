package op

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// DivN 保留固定位数的乘法
func MulN(x float64, y float64, n int) float64 {
	//n小数点位数
	sprintf := fmt.Sprintf("%d", n)
	value, _ := strconv.ParseFloat(fmt.Sprintf("%."+sprintf+"f", x*y), 64)
	return value
}

func Mul(x float64, y float64) float64 {
	s1 := fmt.Sprintf("%v", x)
	n := 0
	_, after, found := strings.Cut(s1, ".")
	if found {
		n = n + len(after)
	}
	s2 := fmt.Sprintf("%v", y)
	_, after, found = strings.Cut(s2, ".")
	if found {
		n = n + len(after)
	}
	//n小数点位数
	sprintf := fmt.Sprintf("%d", n)
	value, _ := strconv.ParseFloat(fmt.Sprintf("%."+sprintf+"f", x*y), 64)
	return value
}

// DivN 保留固定位数的除法
func DivN(x float64, y float64, n int) float64 {
	//n小数点位数
	sprintf := fmt.Sprintf("%d", n)
	value, _ := strconv.ParseFloat(fmt.Sprintf("%."+sprintf+"f", x/y), 64)
	return value
}

// DivN 保留固定位数的加法
func AddN(x float64, y float64, n int) float64 {
	//n小数点位数
	sprintf := fmt.Sprintf("%d", n)
	value, _ := strconv.ParseFloat(fmt.Sprintf("%."+sprintf+"f", x+y), 64)
	return value
}

// FloorFloat 保留precision位小数，后面舍去
func FloorFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Floor(val*ratio) / ratio
}

// RoundFloat 保留precision位小数，后面四舍五入
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
