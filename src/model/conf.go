package model

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type CommandConf struct {
	Role string

	MasterReplid     string
	MasterReplOffset int64

	ReplicaofAddress string
	ReplicaofPort    int
}

func (c CommandConf) String() string {
	buf, _ := json.Marshal(&c)
	return string(buf)
}

func (c CommandConf) ReplicaofAddressAndPort() string {
	if strings.Contains(c.ReplicaofAddress, ":") {
		return fmt.Sprintf("[%s]:%d", c.ReplicaofAddress, c.ReplicaofPort)
	}
	return fmt.Sprintf("%s:%d", c.ReplicaofAddress, c.ReplicaofPort)
}

type pairType struct {
	name  string
	value interface{}
}

func (c CommandConf) Visit(f func(name string, value interface{})) {
	var pairs []*pairType

	pairs = append(pairs, &pairType{"role", c.Role})
	if c.Role == "master" {
		pairs = append(pairs, &pairType{"master_replid", c.MasterReplid})
		pairs = append(pairs, &pairType{"master_repl_offset", c.MasterReplOffset})
	}
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i].name < pairs[j].name
	})

	for _, p := range pairs {
		f(p.name, p.value)
	}
}
