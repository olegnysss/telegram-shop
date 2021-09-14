package couchbase

import "github.com/couchbase/gocb/v2"

type EntityAdapter struct {
}

func (e *EntityAdapter) fetch(document string) (*gocb.GetResult, error) {
	result, err := collection.Get(document, nil)
	if err != nil {
		if kvErr, ok := err.(*gocb.KeyValueError); ok {
			_, err := e.handleKVError(kvErr, document)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return result, err
}

func (e *EntityAdapter) handleKVError(kvErr *gocb.KeyValueError, document string) (mutOut *gocb.MutationResult, errOut error) {
	if kvErr.ErrorDescription == "Not Found" {
		Logger.Printf("%s is not existed. Creating new document.", document)
		result, err := collection.Insert(document, []int{}, nil)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, kvErr
}
