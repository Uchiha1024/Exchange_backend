package bc

import (
	"fmt"
	"testing"
)

func TestWallet_GetAddress(t *testing.T) {

	wallet, err := NewWallet()
	if err != nil {
		panic(err)
	}
	address := wallet.GetTestAddress()
	fmt.Println(string(address))
	priKey := wallet.GetPriKey()
	fmt.Println(priKey)

}
