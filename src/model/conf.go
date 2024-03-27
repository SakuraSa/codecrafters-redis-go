package model

import "encoding/json"

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

func (c CommandConf) Visit(f func(name string, value interface{}) error) error {
	if err := f("role", c.Role); err != nil {
		return err
	}
	if c.Role == "master" {
		if err := f("master_repl_id", c.MasterReplid); err != nil {
			return err
		}
		if err := f("master_repl_offset", c.MasterReplOffset); err != nil {
			return err
		}
	}
	return nil
}
