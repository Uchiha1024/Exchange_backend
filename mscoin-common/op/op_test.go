package op

import (
	"fmt"
	"testing"
)

func TestRoundFloat(t *testing.T) {
	float := FloorFloat(2.0/3.0, 8)
	fmt.Println(float)
}
