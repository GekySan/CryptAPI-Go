package main

import (
	"fmt"
	"examples/cryptapi"
)

// https://docs.cryptapi.io/#operation/estimate

func main() {
	estimate, err := cryptapi.GetEstimate("ltc", 1, "default")
	if err != nil {
		fmt.Println("Error getting estimate:", err)
		return
	}
	for key, value := range estimate {
	    fmt.Printf("%s: %s\n", key, value)
	}
}