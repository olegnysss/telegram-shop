package qiwi

import "time"

type Config struct {
	QiwiToken        string
	QiwiWallet       string
	QiwiPaymentsPath string
	QiwiCashInPath   string
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
			Amount   float64 `json:"amount"`
			Currency float64 `json:"currency"`
		} `json:"sum"`
		Comment string `json:"comment"`
	} `json:"data"`
}
