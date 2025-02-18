package tools

import (
	"fmt"
	"math/rand"
)

// 生成随机四位数，用于生成验证码
func Rand4Num() string {
	num := rand.Intn(9999)
	if num < 1000 {
		num = num + 1000
	}
	return fmt.Sprintf("%s", num)
}
