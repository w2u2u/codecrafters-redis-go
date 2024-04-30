package database

import "errors"

type IDatabase interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

type KeyValue struct {
	data map[string]string
}

func NewKeyValue() KeyValue {
	return KeyValue{
		data: make(map[string]string),
	}
}

func (db *KeyValue) Get(key string) (string, error) {
	value, ok := db.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return value, nil
}

func (db *KeyValue) Set(key string, value string) error {
	db.data[key] = value
	return nil
}
