package command

import "fmt"

func NewSimpleString(msg string) string {
	return fmt.Sprintf("+%s\r\n", msg)
}

func NewBulkString(msg string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)
}

func NewArrays(msgs []string) string {
	resp := fmt.Sprintf("*%d\r\n", len(msgs))

	for _, msg := range msgs {
		resp += NewBulkString(msg)
	}

	return resp
}

func NewNulls() string {
	return "$-1\r\n"
}
