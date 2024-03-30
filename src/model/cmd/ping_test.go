package cmd

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
)

func TestPing_Read(t *testing.T) {
	tests := []struct {
		name    string
		command *Ping
		input   string
		output  string
		isError bool
	}{
		{
			name:    "normal",
			command: &Ping{},
			input:   "*1\r\n+PING\r\n",
			output:  "PING",
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

func TestPing_Execute(t *testing.T) {
	tests := []struct {
		name    string
		command *Ping
		output  string
		isError bool
	}{
		{
			name:    "normal",
			command: &Ping{},
			output:  "+PONG\r\n",
			isError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &strings.Builder{}
			if _, err := tt.command.Execute(writer, nil, nil); err != nil {
				if tt.isError {
					return
				}
				t.Errorf("case %s: failed to execute command: %v", tt.name, err)
			} else if tt.isError {
				t.Errorf("case %s: expected error but got nil", tt.name)
			} else if actual := writer.String(); actual != tt.output {
				t.Errorf("case %s: expected %s but got %s", tt.name, tt.output, actual)
			}
		})
	}
}
