package command

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	Ping = "ping"
	Echo = "echo"
	Get  = "get"
	Set  = "set"
)

const (
	SimpleString = '+'
	BulkString   = '$'
	Arrays       = '*'
)

type RedisCommand struct {
	Args []string
	Size uint
}

func NewRedisCommand(reader *bufio.Reader) (*RedisCommand, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)

	switch line[0] {

	case SimpleString:
		return &RedisCommand{Args: []string{line[1:]}, Size: 1}, nil

	case BulkString:
		arg, err := parseBulkString(reader)
		if err != nil {
			return nil, err
		}
		return &RedisCommand{Args: []string{arg}, Size: 1}, nil

	case Arrays:
		count := 0
		fmt.Sscanf(line, "*%d", &count)
		args := make([]string, count)
		for i := 0; i < count; i++ {
			arg, err := parseBulkString(reader)
			if err != nil {
				return nil, err
			}
			args[i] = arg
		}
		return &RedisCommand{Args: args, Size: uint(count)}, nil
	}

	return &RedisCommand{
		Args: []string{},
	}, nil
}

func parseBulkString(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	line = strings.TrimSpace(line)

	size := 0
	fmt.Sscanf(line, "$%d", &size)
	data := make([]byte, size+2)

	_, err = io.ReadFull(reader, data)
	if err != nil {
		return "", err
	}

	return string(data[:size]), nil
}
