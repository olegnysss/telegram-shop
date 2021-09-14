package couchbase

import (
	"github.com/olegnysss/telebot_qiwi/pkg/qiwi"
	"log"
)

type TransactionsAdapter struct {
	UserTransactionsMap  map[ID]map[ID]Transaction
	entityAdapter        *EntityAdapter
	TransactionsDocument string
}

func InitTransactionsAdapter() *TransactionsAdapter {
	return &TransactionsAdapter{
		map[ID]map[ID]Transaction{},
		&EntityAdapter{},
		"Transactions",
	}
}

func (t *TransactionsAdapter) FetchTransactions(chatId int64) (map[ID]map[ID]Transaction, error) {
	transactionsResult, err := t.entityAdapter.fetch(t.TransactionsDocument)

	if transactionsResult == nil {
		return t.UserTransactionsMap, nil
	}
	var transactions []Transaction
	err = transactionsResult.Content(&transactions)
	if err != nil {
		return nil, err
	}

	transactionsMap := make(map[ID]Transaction)
	for _, transaction := range transactions {
		transactionsMap[transaction.TxnId] = transaction
	}

	t.UserTransactionsMap[ID(chatId)] = transactionsMap

	log.Printf("%+v", t.UserTransactionsMap)
	return t.UserTransactionsMap, nil
}

func (t *TransactionsAdapter) ParseTransactions(telegramId string, responseData qiwi.PaymentsResponse) (map[ID]Transaction, error) {
	transactionMap := make(map[ID]Transaction)
	var sum float64 = 0
	for _, elem := range responseData.Data {
		if elem.Comment == telegramId {
			sum += elem.Sum.Amount
			transactionMap[ID(elem.TxnID)] = Transaction{
				TxnId:    ID(elem.TxnID),
				PersonId: elem.PersonID,
				Sum:      elem.Sum.Amount,
				Comment:  elem.Comment,
			}
		}
	}
	return transactionMap, nil
}

type Transaction struct {
	TxnId    ID
	PersonId int64
	Sum      float64
	Comment  string
}
