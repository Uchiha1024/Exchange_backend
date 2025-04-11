package op

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

func DivN(x float64, y float64, n int) float64 {
	s := fmt.Sprintf("%d", n)
	parseFloat, err := strconv.ParseFloat(fmt.Sprintf("%."+s+"f", x/y), 64)
	if err != nil {
		log.Println(err)
		return 0
	}
	return parseFloat
}

func MulN(x float64, y float64, n int) float64 {
	s := fmt.Sprintf("%d", n)
	parseFloat, err := strconv.ParseFloat(fmt.Sprintf("%."+s+"f", x*y), 64)
	if err != nil {
		log.Println(err)
		return 0
	}
	return parseFloat
}

func AddN(x float64, y float64, n int) float64 {
	s := fmt.Sprintf("%d", n)
	parseFloat, err := strconv.ParseFloat(fmt.Sprintf("%."+s+"f", x+y), 64)
	if err != nil {
		log.Println(err)
		return 0
	}
	return parseFloat
}
func ReduceN(x float64, y float64, n int) float64 {
	s := fmt.Sprintf("%d", n)
	parseFloat, err := strconv.ParseFloat(fmt.Sprintf("%."+s+"f", x-y), 64)
	if err != nil {
		log.Println(err)
		return 0
	}
	return parseFloat
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
	//n小数点位数  x 8 y 8 16位
	sprintf := fmt.Sprintf("%d", n)
	value, _ := strconv.ParseFloat(fmt.Sprintf("%."+sprintf+"f", x*y), 64)
	return value
}

func MulFloor(x float64, y float64, n int) float64 {
	//自动根据小数点的位数 进行保留
	mulN := Mul(x, y)
	return FloorFloat(mulN, uint(n))
}

func Div(x float64, y float64) float64 {
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
	//n小数点位数  x 8 y 8 16位
	sprintf := fmt.Sprintf("%d", n)
	value, _ := strconv.ParseFloat(fmt.Sprintf("%."+sprintf+"f", x/y), 64)
	return value
}

func DivFloor(x float64, y float64, n int) float64 {
	//自动根据小数点的位数 进行保留
	mulN := Div(x, y)
	return FloorFloat(mulN, uint(n))
}
func Sub(x float64, y float64) float64 {
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
	//n小数点位数  x 8 y 8 16位
	sprintf := fmt.Sprintf("%d", n)
	value, _ := strconv.ParseFloat(fmt.Sprintf("%."+sprintf+"f", x-y), 64)
	return value
}

func SubFloor(x float64, y float64, n int) float64 {
	//自动根据小数点的位数 进行保留
	mulN := Sub(x, y)
	return FloorFloat(mulN, uint(n))
}
func Add(x float64, y float64) float64 {
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
	//n小数点位数  x 8 y 8 16位
	sprintf := fmt.Sprintf("%d", n)
	value, _ := strconv.ParseFloat(fmt.Sprintf("%."+sprintf+"f", x+y), 64)
	return value
}

func AddFloor(x float64, y float64, n int) float64 {
	//自动根据小数点的位数 进行保留
	mulN := Add(x, y)
	return FloorFloat(mulN, uint(n))
}

// FloorFloat 1.245  保留两位小数点 1.25  金融方面的应用 1.24  2.0/3.0
func FloorFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Floor(val*ratio) / ratio
}

// RoundFloat 进位运算  1.245  保留两位小数点 1.25  金融方面的应用 1.24 而不是 1.25
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
