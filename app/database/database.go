package database

import "encoding/hex"

type IDatabase interface {
	Get(key string) (string, error)
	Set(key string, value string, exp string)
}

func EmptyRDB() []byte {
	data, _ := hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
	return data
}
