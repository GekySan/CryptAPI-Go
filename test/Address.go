package main

import (
	"fmt"
	"examples/cryptapi"
)

// https://docs.cryptapi.io/#operation/create

func main() {
	parameters := map[string]string{}
	caParams := map[string]string{}

	helper, err := cryptapi.NewCryptAPIHelper("ltc", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "http://a.b.c.d:2468/cryptapi/5961572520/1", parameters, caParams)
	if err != nil {
		fmt.Println("Error creating CryptAPIHelper:", err)
		return
	}
	
	address, err := helper.GetAddress()
	if err != nil {
		fmt.Println("Error getting address:", err)
		return
	}
	fmt.Println("Payment Address:", address["address_in"])
}