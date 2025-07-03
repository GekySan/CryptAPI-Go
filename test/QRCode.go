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
	
	qrCode, err := helper.GetQRCode("", 600)
	if err != nil {
		fmt.Println("Error getting QR code:", err)
		return
	}
	fmt.Println("QR Code:", qrCode["qr_code"])
}
