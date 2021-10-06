package couchbase

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/olegnysss/telebot_qiwi/pkg/qiwi"
	"strconv"
)

var tnxCollection *gocb.Collection

type TransactionsAdapter struct {
	entityAdapter          *EntityAdapter
	TransactionsCollection string
}

func InitTransactionsAdapter() *TransactionsAdapter {
	return &TransactionsAdapter{
		&EntityAdapter{},
		"transactions",
	}
}

func InitTnxCollection(scope *gocb.Scope) *gocb.Collection {
	tnxCollection = scope.Collection("transactions")
	return tnxCollection
}

func (t *TransactionsAdapter) ProcessTnx(transactions map[ID]Transaction, couch *CouchClient) (bool, error) {
	var max ID = 0
	for _, transaction := range transactions {
		if transaction.TxnId > max {
			max = transaction.TxnId
		}
	}
	isNew, err := t.tnxFlow(transactions[max], couch)
	return isNew, err
}

func (t *TransactionsAdapter) tnxFlow(tnx Transaction, couch *CouchClient) (bool, error) {
	_, isNew, err := t.CheckTnx(int64(tnx.TxnId))
	if err != nil {
		return false, err
	}
	if isNew {
		userResult, err := couch.UsersAdapter.retrieve(int64(tnx.UserId))
		if err != nil {
			return false, err
		}
		user, err := ContentUser(userResult)
		user.Balance += tnx.Sum
		_, err = couch.UsersAdapter.entityAdapter.replace(usersCollection, user)
		if err != nil {
			return false, err
		}
		_, err = t.entityAdapter.insert(tnxCollection, tnx)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (t *TransactionsAdapter) CheckTnx(tnxId int64) (Transaction, bool, error) {
	tnxResult, err := t.entityAdapter.retrieve(tnxCollection, int(tnxId))
	if err != nil {
		if kvErr, ok := err.(*gocb.KeyValueError); ok && kvErr.ErrorDescription == "Not Found" {
			return Transaction{}, true, nil
		} else {
			return Transaction{}, false, err
		}
	}
	return contentTnx(tnxResult)
}

func contentTnx(tnxResult *gocb.GetResult) (Transaction, bool, error) {
	var tnx Transaction
	err := tnxResult.Content(&tnx)
	if err != nil {
		panic(err)
	}
	return tnx, false, nil
}

func (t *TransactionsAdapter) FetchTransactions(id int64) ([]Transaction, error) {
	query := fmt.Sprintf("SELECT transactions.* from `teleshopBucket`._default.transactions WHERE UserId = %d", id)
	results, err := cluster.Query(query, nil)
	if err != nil {
		return nil, err
	}
	return contentTnxQuery(results)
}

func contentTnxQuery(tnxResults *gocb.QueryResult) ([]Transaction, error) {
	var tnxSlice []Transaction
	for tnxResults.Next() {
		var tnx Transaction
		err := tnxResults.Row(&tnx)
		if err != nil {
			panic(err)
		}
		tnxSlice = append(tnxSlice, tnx)
	}
	err := tnxResults.Err()
	if err != nil {
		return nil, err
	}
	return tnxSlice, err
}

func (t *TransactionsAdapter) ParseTransactions(telegramId string, responseData qiwi.PaymentsResponse) (map[ID]Transaction, error) {
	transactionMap := make(map[ID]Transaction)
	userId, err := strconv.Atoi(telegramId)
	if err != nil {
		return nil, err
	}
	for _, elem := range responseData.Data {
		if elem.Comment == telegramId {
			transactionMap[ID(elem.TxnID)] = Transaction{
				TxnId:  ID(elem.TxnID),
				Msisdn: elem.PersonID,
				Sum:    elem.Sum.Amount,
				UserId: ID(userId),
			}
		}
	}
	return transactionMap, nil
}

type Transaction struct {
	TxnId  ID
	Msisdn int64
	Sum    float64
	UserId ID
}

func (t Transaction) GetId() ID {
	return t.TxnId
}

func (t Transaction) GetStringId() string {
	return strconv.Itoa(int(t.TxnId))
}
