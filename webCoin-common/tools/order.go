package tools

import (
	"fmt"
	"math/rand"
	"time"
)

// Unq 根据时间戳和随机数生成唯一id
func Unq(prefix string) string {
	milli := time.Now().UnixMilli()
	intn := rand.Intn(999999)
	return fmt.Sprintf("%s%d%d", prefix, milli, intn)
}
