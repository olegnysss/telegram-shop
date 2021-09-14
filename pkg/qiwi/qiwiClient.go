package qiwi

import (
	"encoding/json"
	"fmt"
	"github.com/olegnysss/telebot_qiwi/pkg/config"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//Struct for payment (id, sum)
//Not more than 100 requests in minute (once in a minute per acc)

var Logger *log.Logger

type QiwiClient struct {
	Token        string
	Wallet       string
	PaymentsPath string
	CashInPath   string
}

func InitQiwiClient(config config.Qiwi) *QiwiClient {
	return &QiwiClient{
		Token:        config.QiwiToken,
		Wallet:       config.QiwiWallet,
		PaymentsPath: config.QiwiPaymentsPath,
		CashInPath:   config.QiwiCashInPath,
	}
}

func (q *QiwiClient) CheckPayment() (PaymentsResponse, error) {
	err := initLogs()
	if err != nil {
		return PaymentsResponse{}, err
	}

	client := &http.Client{}
	req, err := q.prepareRequest()
	if err != nil {
		return PaymentsResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return PaymentsResponse{}, err
	}

	defer resp.Body.Close()
	var responseData PaymentsResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return PaymentsResponse{}, err
	}

	if err := json.Unmarshal(body, &responseData); err != nil {
		return PaymentsResponse{}, err
	}
	return responseData, nil
}

func (q *QiwiClient) prepareRequest() (*http.Request, error) {
	paymentsURL := fmt.Sprintf(q.PaymentsPath, q.Wallet)
	req, err := http.NewRequest("GET", paymentsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer  %s", q.Token))
	req.Header.Add("Content-Type", "application/json")
	Logger.Printf("request : %+v", req.URL)
	return req, nil
}

func initLogs() error {
	file, err := os.OpenFile("qiwi.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	Logger = log.New(file, "QIWI: ", log.Ldate|log.Ltime|log.Lshortfile)
	return err
}
