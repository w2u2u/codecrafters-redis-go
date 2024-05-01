package command

import "bufio"

func IsOk(reader *bufio.Reader) bool {
	data, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	return data == "+OK\r\n"
}

func IsPong(reader *bufio.Reader) bool {
	data, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	return data == "+PONG\r\n"
}
