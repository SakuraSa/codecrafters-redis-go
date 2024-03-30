package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/src/model"
	"github.com/codecrafters-io/redis-starter-go/src/model/redis"
)

func TestSet_Read(t *testing.T) {
	tests := []struct {
		name    string
		command *Set
		input   string
		output  string
		isError bool
	}{
		{
			name:    "normal",
			command: &Set{},
			input:   "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n",
			output:  "SET[key, value, 9223372036854775807]",
		},
		{
			name:    "set with px option",
			command: &Set{},
			input:   "*5\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n$2\r\nPX\r\n$3\r\n100\r\n",
			output:  fmt.Sprintf("SET[key, value, %d]", 100+time.Now().UnixNano()/int64(time.Millisecond)),
		},
		{
			name:    "wrong number of arguments",
			command: &Set{},
			input:   "*2\r\n$3\r\nSET\r\n$3\r\nkey\r\n",
			isError: true,
		},
		{
			name:    "wrong px option",
			command: &Set{},
			input:   "*5\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n$2\r\nPX\r\n$3\r\nabc\r\n",
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

func TestSet_Execute(t *testing.T) {
	tests := []struct {
		name       string
		command    *Set
		storage    *model.RedisStorage
		conf       *model.CommandConf
		output     string
		memChecker func(*model.RedisStorage) error
		isError    bool
	}{
		{
			name:    "normal",
			command: &Set{key: "key", value: []byte("value"), expireAt: 123123},
			storage: &model.RedisStorage{Mem: make(map[string]*model.RedisBucket)},
			conf:    &model.CommandConf{},
			output:  "+OK\r\n",
			memChecker: func(storage *model.RedisStorage) error {
				if bucket, ok := storage.Mem["key"]; !ok {
					return fmt.Errorf("key not found")
				} else if string(bucket.Value) != "value" {
					return fmt.Errorf("value not match")
				} else if bucket.ExpireAt != 123123 {
					return fmt.Errorf("expireAt not match")
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
