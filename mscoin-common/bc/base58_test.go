package bc

import (
	"fmt"
	"testing"
)

func TestBase58Encode(t *testing.T) {
	encode := Base58Encode([]byte("asdasdasdasdasdasdas"))
	fmt.Println(string(encode))
	decode := Base58Decode(encode)
	fmt.Println(string(decode))
}
