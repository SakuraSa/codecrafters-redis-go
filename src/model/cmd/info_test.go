package cmd

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/src/model"
	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
	"github.com/codecrafters-io/redis-starter-go/src/util"
)

func TestInfo_Read(t *testing.T) {
	tests := []struct {
		name    string
		command *Info
		input   string
		output  string
		isError bool
	}{
		{
			name:    "normal",
			command: &Info{},
			input:   "*2\r\n+INFO\r\n+replication\r\n",
			output:  "INFO[replication]",
		},
		{
			name:    "wrong number of arguments",
			command: &Info{},
			input:   "*1\r\n$4\r\nINFO\r\n",
			isError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewBufferString(tt.input))
			if obj, err := redis.ReadObject(reader, redis.ArrayLeading); err != nil {
				t.Errorf("case %s: failed to read object: %v", tt.name, err)
			} else if args, ok := obj.(*redis.Array); !ok {
				t.Errorf("case %s: expected *redis.Array but got %v", tt.name, obj)
			} else if err := tt.command.Read(args); err != nil {
				if tt.isError {
					return
				}
				t.Errorf("case %s: failed to read command: %v", tt.name, err)
			} else if tt.isError {
				t.Errorf("case %s: expected error but got nil", tt.name)
			} else if actual := tt.command.String(); actual != tt.output {
				t.Errorf("case %s: expected %s but got %s", tt.name, tt.output, actual)
			}
		})
	}
}

func TestInfo_Execute(t *testing.T) {
	tests := []struct {
		name    string
		command *Info
		conf    *model.CommandConf
		output  string
		isError bool
	}{
		{
			name:    "slave replication",
			command: &Info{subCommand: &InfoReplication{}},
			conf: &model.CommandConf{
				Role:             "slave",
				MasterReplid:     "id",
				MasterReplOffset: 123,
				ReplicaofAddress: "localhost",
				ReplicaofPort:    456,
			},
			output:  "$12\r\nrole:slave\r\n\r\n",
			isError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &strings.Builder{}
			if _, err := tt.command.Execute(writer, nil, tt.conf); err != nil {
				if tt.isError {
					return
				}
				t.Errorf("case %s: failed to execute command: %v", tt.name, err)
			} else if tt.isError {
				t.Errorf("case %s: expected error but got nil", tt.name)
			} else if actual := writer.String(); actual != tt.output {
				t.Errorf("case %s: expected %v but got %v", tt.name, util.J(tt.output), util.J(actual))
			}
		})
	}
}
