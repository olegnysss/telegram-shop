package qiwi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//Struct for payment (id, sum)
//Not more than 100 requests in minute (once in a minute per acc)

func CheckPayment(config Config) {
	client := &http.Client{}
	paymentsURL := fmt.Sprintf(config.QiwiPaymentsPath, config.QiwiWallet)
	req, err := http.NewRequest("GET", paymentsURL, nil)
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer  %s", config.QiwiToken))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	var responseData PaymentsResponseStruct
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", responseData)
	telegramId := "460158421"
	sum := 0
	for _, elem := range responseData.Data {
		if elem.Comment == telegramId {
			sum += elem.Sum.Amount
		}
	}
	fmt.Printf("Сумма для аккаунта %s: %d", telegramId, sum)
}