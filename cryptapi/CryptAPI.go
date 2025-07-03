package cryptapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type CryptAPIException struct {
	Message string
}

func (e *CryptAPIException) Error() string {
	return e.Message
}

type CryptAPIHelper struct {
	Coin          string
	OwnAddress    string
	CallbackURL   string
	Parameters    map[string]string
	CAParams      map[string]string
	PaymentAddress string
}

const (
	CRYPTAPI_URL = "https://api.cryptapi.io/"
	CRYPTAPI_HOST = "api.cryptapi.io"
)

func NewCryptAPIHelper(coin, ownAddress, callbackURL string, parameters, caParams map[string]string) (*CryptAPIHelper, error) {
	if callbackURL == "" {
		return nil, errors.New("Callback URL is missing")
	}
	if coin == "" {
		return nil, errors.New("Coin is missing")
	}
	if ownAddress == "" {
		return nil, errors.New("Address is missing")
	}

	return &CryptAPIHelper{
		Coin:          coin,
		OwnAddress:    ownAddress,
		CallbackURL:   callbackURL,
		Parameters:    parameters,
		CAParams:      caParams,
		PaymentAddress: "",
	}, nil
}

func (h *CryptAPIHelper) GetAddress() (map[string]interface{}, error) {
	coin := h.Coin

	params := url.Values{}
	params.Add("address", h.OwnAddress)
	params.Add("callback", h.CallbackURL)

	for key, value := range h.Parameters {
		params.Add(key, value)
	}
	for key, value := range h.CAParams {
		params.Add(key, value)
	}

	address, err := processRequest(coin, "create", params)
	if err != nil {
		return nil, err
	}

	if address != nil {
		h.PaymentAddress = address["address_in"].(string)
	}

	return address, nil
}

func (h *CryptAPIHelper) GetLogs() (map[string]interface{}, error) {
	coin := h.Coin

	params := url.Values{}
	params.Add("callback", h.CallbackURL)

	for key, value := range h.Parameters {
		params.Add(key, value)
	}

	return processRequest(coin, "logs", params)
}

func (h *CryptAPIHelper) GetQRCode(value string, size int) (map[string]interface{}, error) {
	params := url.Values{}
	params.Add("address", h.PaymentAddress)
	params.Add("size", fmt.Sprintf("%d", size))

	if value != "" {
		params.Add("value", value)
	}

	return processRequest(h.Coin, "qrcode", params)
}

func (h *CryptAPIHelper) GetConversion(fromCoin string, value string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Add("from", fromCoin)
	params.Add("value", value)

	return processRequest(h.Coin, "convert", params)
}

func GetInfo(coin string) (map[string]interface{}, error) {
	return processRequest(coin, "info", nil)
}

func GetSupportedCoins() (map[string]string, error) {
	info, err := GetInfo("")
	if err != nil {
		return nil, err
	}

	delete(info, "fee_tiers")

	coins := make(map[string]string)
	for ticker, coinInfo := range info {
		if coin, ok := coinInfo.(map[string]interface{})["coin"]; ok {
			coins[ticker] = coin.(string)
		} else {
			for token, tokenInfo := range coinInfo.(map[string]interface{}) {
				tokenCoin := tokenInfo.(map[string]interface{})["coin"].(string)
				coins[ticker+"_"+token] = tokenCoin + " (" + strings.ToUpper(ticker) + ")"
			}
		}
	}

	return coins, nil
}

func GetEstimate(coin string, addresses int, priority string) (map[string]interface{}, error) {
	params := url.Values{}
	params.Add("addresses", fmt.Sprintf("%d", addresses))
	params.Add("priority", priority)

	return processRequest(coin, "estimate", params)
}

func processRequest(coin, endpoint string, params url.Values) (map[string]interface{}, error) {
	if coin != "" {
		coin += "/"
	}

	reqURL := fmt.Sprintf("%s%s%s/", CRYPTAPI_URL, strings.Replace(coin, "_", "/", -1), endpoint)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = params.Encode()
	req.Header.Add("Host", CRYPTAPI_HOST)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseObj map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseObj); err != nil {
		return nil, err
	}

	if status, ok := responseObj["status"]; ok && status == "error" {
		return nil, &CryptAPIException{Message: responseObj["error"].(string)}
	}

	return responseObj, nil
}