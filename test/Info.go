package main

import (
	"fmt"
    "strconv"
	"examples/cryptapi"
)

// https://docs.cryptapi.io/#operation/info

func main() {
	info, err := cryptapi.GetInfo("ltc")
	if err != nil {
		fmt.Println("Error getting info:", err)
		return
	}
	// fmt.Println(conversion)

    euroPrice := strconv.FormatFloat(func() float64 {
        prices := info["prices"].(map[string]interface{})
        euroPriceStr := prices["EUR"].(string)
        euroPrice, _ := strconv.ParseFloat(euroPriceStr, 64)
        return euroPrice
    }(), 'f', 2, 64)

    fmt.Println("LTC is actually " + euroPrice + " EUR")
}
