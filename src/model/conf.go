package model

import "encoding/json"

type CommandConf struct {
	Role             string
	ReplicaofAddress string
	ReplicaofPort    int
}

func (c CommandConf) String() string {
	buf, _ := json.Marshal(&c)
	return string(buf)
}
