package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/src/model"
	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
)

func TestGet_Read(t *testing.T) {
	tests := []struct {
		name    string
		command *Get
		input   string
		output  string
		isError bool
	}{
		{
			name:    "normal",
			command: &Get{},
			input:   "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n",
			output:  "GET[key]",
		},
		{
			name:    "wrong number of arguments",
			command: &Get{},
			input:   "*1\r\n$3\r\nGET\r\n$3\r\n",
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

func TestGet_Execute(t *testing.T) {
	tests := []struct {
		name       string
		command    *Get
		storage    *model.RedisStorage
		conf       *model.CommandConf
		output     string
		memChecker func(*model.RedisStorage) error
		isError    bool
	}{
		{
			name:       "normal",
			command:    &Get{key: "key"},
			storage:    &model.RedisStorage{Mem: map[string]*model.RedisBucket{"key": {Value: []byte("value"), ExpireAt: math.MaxInt64}}},
			conf:       &model.CommandConf{},
			output:     "$5\r\nvalue\r\n",
			memChecker: nil,
			isError:    false,
		},
		{
			name:       "missing key",
			command:    &Get{key: "key"},
			storage:    &model.RedisStorage{Mem: map[string]*model.RedisBucket{}},
			conf:       &model.CommandConf{},
			output:     "$-1\r\n",
			memChecker: nil,
			isError:    false,
		},
		{
			name:    "expire key",
			command: &Get{key: "key"},
			storage: &model.RedisStorage{Mem: map[string]*model.RedisBucket{"key": {Value: []byte("value"), ExpireAt: 1}}},
			conf:    &model.CommandConf{},
			output:  "$-1\r\n",
			memChecker: func(rs *model.RedisStorage) error {
				if _, ok := rs.Mem["key"]; ok {
					return fmt.Errorf("expected key to be deleted but it still exists")
				}
				return nil
			},
			isError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &strings.Builder{}
			if _, err := tt.command.Execute(writer, tt.storage, tt.conf); err != nil {
				if tt.isError {
					return
				}
				t.Errorf("case %s: failed to execute command: %v", tt.name, err)
			} else if tt.isError {
				t.Errorf("case %s: expected error but got nil", tt.name)
			} else if actual := writer.String(); actual != tt.output {
				t.Errorf("case %s: expected %s but got %s", tt.name, tt.output, actual)
			} else if tt.memChecker != nil {
				if err := tt.memChecker(tt.storage); err != nil {
					t.Errorf("case %s: %v", tt.name, err)
				}
			}
		})
	}
}
