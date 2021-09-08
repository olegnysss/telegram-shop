package qiwi

import "time"

type Config struct {
	QiwiToken        string
	QiwiWallet       string
	QiwiPaymentsPath string
}

type PaymentsResponseStruct struct {
	Data []struct {
		TxnID    int64     `json:"txnId"`
		PersonID int64     `json:"personId"`
		Date     time.Time `json:"date"`
		Status   string    `json:"status"`
		Type     string    `json:"type"`
		TrmTxnID string    `json:"trmTxnId"`
		Account  string    `json:"account"`
		Sum      struct {
			Amount   int `json:"amount"`
			Currency int `json:"currency"`
		} `json:"sum"`
		Comment string `json:"comment"`
	} `json:"data"`
}
