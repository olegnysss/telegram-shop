package qiwi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//Struct for payment (id, sum)
//Not more than 100 requests in minute (once in a minute per acc)

var Logger *log.Logger

func CheckPayment(config Config) {
	file, err := os.OpenFile("qiwi.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	Logger = log.New(file, "QIWI: ", log.Ldate|log.Ltime|log.Lshortfile)

	client := &http.Client{}
	paymentsURL := fmt.Sprintf(config.QiwiPaymentsPath, config.QiwiWallet)
	req, err := http.NewRequest("GET", paymentsURL, nil)
	if err != nil {
		Logger.Panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer  %s", config.QiwiToken))
	req.Header.Add("Content-Type", "application/json")
	//todo Bearer token here, vulnerable
	Logger.Printf("request : %+v", req)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	var responseData PaymentsResponseStruct
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Println(err)
		return
	}

	Logger.Printf("%+v\n", responseData)
	//todo should receive user
	telegramId := "460158421"
	var sum float64 = 0
	for _, elem := range responseData.Data {
		if elem.Comment == telegramId {
			sum += elem.Sum.Amount
		}
	}
	Logger.Printf("Сумма для аккаунта %s: %f", telegramId, sum)
}
