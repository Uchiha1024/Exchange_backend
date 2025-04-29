package tools

import (
	"fmt"
	"math/rand"
	"time"
)

// 生成订单号   

func Unq(prefix string) string {
	milli := time.Now().UnixMilli()
	intn := rand.Intn(999999)
	return fmt.Sprintf("%s%d%d", prefix, milli, intn)
}
