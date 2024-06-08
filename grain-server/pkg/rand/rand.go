package randx

import (
	"github.com/go-grain/grain/pkg/convert"
	"math/rand"
	"time"
)

// RandomString 获取随机数
func RandomString(length int) string {
	rand.NewSource(time.Now().UnixNano())
	const charset = "0123456789"
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, charset[rand.Intn(len(charset))])
	}
	return string(result)
}

func RandomInt64(length int) int64 {
	rand.NewSource(time.Now().UnixNano())
	const charset = "0123456789"
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, charset[rand.Intn(len(charset))])
	}
	return convert.String2Int64(string(result))
}

func RandomCharset(length int) string {
	rand.NewSource(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var result []byte
	for i := 0; i < length; i++ {
		result = append(result, charset[rand.Intn(len(charset))])
	}
	return string(result)
}
