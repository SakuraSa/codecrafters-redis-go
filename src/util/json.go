package util

import (
	"encoding/json"
	"fmt"
)

var (
	_ fmt.Stringer = &JsonWrapper{}
)

type JsonWrapper struct {
	value interface{}
}

func (j *JsonWrapper) String() string {
	buf, _ := json.Marshal(j.value)
	return string(buf)
}

func J(value interface{}) *JsonWrapper {
	return &JsonWrapper{value: value}
}
