package main

import (
	"fmt"
	"examples/cryptapi"
)

// https://docs.cryptapi.io/#operation/convert

func main() {
	parameters := map[string]string{}
	caParams := map[string]string{}

	helper, err := cryptapi.NewCryptAPIHelper("ltc", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "http://a.b.c.d:2468/cryptapi/5961572520/1", parameters, caParams)
	if err != nil {
		fmt.Println("Error creating CryptAPIHelper:", err)
		return
	}

	conversion, err := helper.GetConversion("ltc", "1")
	if err != nil {
		fmt.Println("Error converting currency:", err)
		return
	}
	// fmt.Println(conversion)
	for key, value := range conversion {
        switch v := value.(type) {
        case string:
            fmt.Printf("%s: %s\n", key, v)
        case int:
            fmt.Printf("%s: %d\n", key, v)
        case float64:
            fmt.Printf("%s: %f\n", key, v)
        default:
            fmt.Printf("%s: %v\n", key, v)
        }
    }
}