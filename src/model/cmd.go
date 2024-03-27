package model

import (
	"encoding/json"
)

type RedisCommandAndArgs struct {
	Args [][]byte
}

func (c RedisCommandAndArgs) String() string {
	var arr []string
	for _, arg := range c.Args {
		arr = append(arr, string(arg))
	}
	buf, _ := json.Marshal(arr)
	return "CMD" + string(buf)
}

func (c RedisCommandAndArgs) Command() string {
	if len(c.Args) == 0 {
		return ""
	}
	return string(c.Args[0])
}
