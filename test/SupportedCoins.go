package main

import (
	"fmt"
	"examples/cryptapi"
)

func main() {
	supportedCoins, err := cryptapi.GetSupportedCoins()
	if err != nil {
		fmt.Println("Error getting supported coins:", err)
		return
	}
	for key, value := range supportedCoins {
	    fmt.Printf("%s: %s\n", key, value)
	}
}
