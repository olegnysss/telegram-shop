package couchbase

import (
	"github.com/couchbase/gocb/v2"
	"strconv"
)

type EntityAdapter struct {
}

type CouchEntity interface {
	GetId() ID
	GetStringId() string
}

func (e *EntityAdapter) retrieve(collection *gocb.Collection, key int) (*gocb.GetResult, error) {
	stringId := strconv.Itoa(key)
	result, err := collection.Get(stringId, nil)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (e *EntityAdapter) insert(collection *gocb.Collection, entity CouchEntity) (*gocb.MutationResult, error) {
	//inapplicable for filling balance
	Logger.Printf("Creating new document with Id %d with type %T", entity.GetId(), entity)
	result, err := collection.Insert(entity.GetStringId(), entity, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *EntityAdapter) replace(collection *gocb.Collection, entity CouchEntity) (*gocb.MutationResult, error) {
	Logger.Printf("Updating document with Id %d with type %T", entity.GetId(), entity)
	result, err := collection.Replace(entity.GetStringId(), entity, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}
