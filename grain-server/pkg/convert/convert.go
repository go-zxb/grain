package convert

import (
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func String2Int(intStr string) (intNum int) {
	intNum, _ = strconv.Atoi(intStr)
	return
}

func String2Int64(intStr string) (int64Num int64) {
	intNum, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0
	}
	int64Num = int64(intNum)
	return
}

func String2Float64(floatStr string) (floatNum float64) {
	floatNum, _ = strconv.ParseFloat(floatStr, 64)
	return
}

func String2Float32(floatStr string) (floatNum float32) {
	floatNum64, _ := strconv.ParseFloat(floatStr, 32)
	floatNum = float32(floatNum64)
	return
}

func Int2String(intNum int) (intStr string) {
	intStr = strconv.Itoa(intNum)
	return
}

func Int642String(intNum int64) (int64Str string) {
	//10, 代表10进制
	int64Str = strconv.FormatInt(intNum, 10)
	return
}

func Float64ToString(floatNum float64, prec ...int) (floatStr string) {
	if len(prec) > 0 {
		floatStr = strconv.FormatFloat(floatNum, 'f', prec[0], 64)
		return
	}
	floatStr = strconv.FormatFloat(floatNum, 'f', -1, 64)
	return
}

func Float32ToString(floatNum float32, prec ...int) (floatStr string) {
	if len(prec) > 0 {
		floatStr = strconv.FormatFloat(float64(floatNum), 'f', prec[0], 32)
		return
	}
	floatStr = strconv.FormatFloat(float64(floatNum), 'f', -1, 32)
	return
}

func BinaryToDecimal(bit string) (num int) {
	fields := strings.Split(bit, "")
	lens := len(fields)
	var tempF float64 = 0
	for i := 0; i < lens; i++ {
		floatNum := String2Float64(fields[i])
		tempF += floatNum * math.Pow(2, float64(lens-i-1))
	}
	num = int(tempF)
	return
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}
