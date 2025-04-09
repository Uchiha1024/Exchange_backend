package tools

import (
	"fmt"
	"math/rand"
)

func Rand4Code() string {
	randInt := rand.Intn(9999)
	if randInt < 1000 {
		randInt += 1000
	}
	return fmt.Sprintf("%04d", randInt)
}
