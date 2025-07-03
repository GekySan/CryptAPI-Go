package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

/*
	publicKeyPEM: https://api.cryptapi.io/pubkey/
	allowedIPs: Note of https://docs.cryptapi.io/#tag/Callbacks

	Configure the callback URL in the function parameters : IP:PORT/cryptapi/{UserID}/{OrderID}. Adjust the code, currently only logs the events.
*/

const (
	publicKeyPEM = "-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC3FT0Ym8b3myVxhQW7ESuuu6lo\ndGAsUJs4fq+Ey//jm27jQ7HHHDmP1YJO7XE7Jf/0DTEJgcw4EZhJFVwsk6d3+4fy\nBsn0tKeyGMiaE6cVkX0cy6Y85o8zgc/CwZKc0uw6d5siAo++xl2zl+RGMXCELQVE\nox7pp208zTvown577wIDAQAB\n-----END PUBLIC KEY-----"
	logFileName = "logs.txt"
)

var (
	allowedIPs = [2]string{"145.239.119.223", "135.125.112.47"}
)

func loadPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("FAIL")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func verifySignature(data, signature string, pubKey *rsa.PublicKey) bool {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed, sig)
	return err == nil
}

func logRequest(w http.ResponseWriter, r *http.Request, pubKey *rsa.PublicKey, id string, orderNum string) {
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Can't open logs file", http.StatusInternalServerError)
		return
	}
	defer logFile.Close()

	ip := strings.Split(r.RemoteAddr, ":")[0]
	signature := r.Header.Get("x-ca-signature")
	var data string
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		data = string(body)
	} else {
		data = r.URL.String()
	}

	logData := fmt.Sprintf("IP: %s\nSignature: %s\nData: %s\nID: %s\nOrder Number: %s\n\n", ip, signature, data, id, orderNum)
	logFile.WriteString(logData)

	w.Header().Set("Content-Type", "application/json")
	if verifySignature(data, signature, pubKey) {
		w.Write([]byte(`{"status": "success"}`))
	} else {
		w.Write([]byte(`{"status": "error", "message": "Invalid signature"}`))
	}
}

func callbackHandler(pubKey *rsa.PublicKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		allowed := false
		for _, allowedIP := range allowedIPs {
			if ip == allowedIP {
				allowed = true
				break
			}
		}
		if !allowed {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status": "forbidden"}`))
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		// fmt.Println(parts)
		if len(parts) < 4 {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}
		id := parts[2]
		orderNum := parts[3]

		logRequest(w, r, pubKey, id, orderNum)
	}
}

func main() {
	pubKey, err := loadPublicKey(publicKeyPEM)
	if err != nil {
		panic(fmt.Sprintf("Can't to load public key: %v", err))
	}

	http.HandleFunc("/cryptapi/", callbackHandler(pubKey))
	http.ListenAndServe(":2468", nil)
}
