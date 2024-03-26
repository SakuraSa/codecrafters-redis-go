package model

import "strings"

type RedisCommandAndArgs struct {
	Args [][]byte
}

func (c RedisCommandAndArgs) String() string {
	var builder strings.Builder
	for _, arg := range c.Args {
		builder.Write(arg)
		builder.Write([]byte(" "))
	}
	return builder.String()
}

func (c RedisCommandAndArgs) Command() string {
	if len(c.Args) == 0 {
		return ""
	}
	return string(c.Args[0])
}
