package database

type IDatabase interface {
	Get(key string) (string, error)
	Set(key string, value string, exp string)
}
