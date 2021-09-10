package couchbase

import "log"

var transactionsDocument = "Transactions"
var transactionsMap map[ID][]Transaction

func FetchTransactions(chatId int64) ([]Transaction, error) {
	transactionsResult, err := collection.Get(transactionsDocument, nil)
	if err != nil {
		return nil, err
	}

	var transactions []Transaction
	err = transactionsResult.Content(&transactions)
	if err != nil {
		return nil, err
	}

	transactionsMap[ID(chatId)] = transactions

	log.Printf("%+v", usersMap)
	return transactions, nil
}
