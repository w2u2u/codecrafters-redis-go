package database

import (
	"errors"
	"time"
)

type dataValue struct {
	value      string
	expiration time.Time
}

func (d dataValue) isExpired() bool {
	return d.expiration != time.Time{} && time.Now().After(d.expiration)
}

type KeyValue struct {
	data map[string]dataValue
}

func NewKeyValue() KeyValue {
	return KeyValue{
		data: make(map[string]dataValue),
	}
}

func (db *KeyValue) Get(key string) (string, error) {
	if data, ok := db.data[key]; ok {
		if !data.isExpired() {
			return data.value, nil
		}
		delete(db.data, key)
	}

	return "", errors.New("Key not found or expired")
}

func (db *KeyValue) Set(key string, value string, exp string) {
	expiration := time.Time{}
	if exp != "0" {
		if duration, err := time.ParseDuration(exp); err == nil {
			expiration = time.Now().Add(duration)
		}
	}

	db.data[key] = dataValue{
		value,
		expiration,
	}
}
