package redis

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestArray_Read(t *testing.T) {
	tests := []struct {
		name     string
		a        *Array
		line     []byte
		expected string
		isError  bool
	}{
		{
			name:     "empty array",
			a:        &Array{},
			line:     []byte("*0\r\n"),
			expected: "Array[]",
			isError:  false,
		},
		{
			name:    "malformed array size",
			a:       &Array{},
			line:    []byte("*a\r\n"),
			isError: true,
		},
		{
			name:    "array size insufficient",
			a:       &Array{},
			line:    []byte("*1\r\n"),
			isError: true,
		},
		{
			name:     "array with simplestring",
			a:        &Array{},
			line:     []byte("*1\r\n+hello\r\n"),
			expected: "Array[SimpleString{hello}]",
			isError:  false,
		},
		{
			name:     "array with bulkstring",
			a:        &Array{},
			line:     []byte("*1\r\n$5\r\nhello\r\n"),
			expected: "Array[BulkString{hello}]",
			isError:  false,
		},
		{
			name:     "array with integer",
			a:        &Array{},
			line:     []byte("*1\r\n:123\r\n"),
			expected: "Array[Integer{123}]",
			isError:  false,
		},
		{
			name:     "array with null",
			a:        &Array{},
			line:     []byte("*1\r\n_\r\n"),
			expected: "Array[Null{}]",
			isError:  false,
		},
		{
			name:     "array with boolean",
			a:        &Array{},
			line:     []byte("*2\r\n#t\r\n#f\r\n"),
			expected: "Array[Boolean{true}, Boolean{false}]",
			isError:  false,
		},
		{
			name:     "array with double",
			a:        &Array{},
			line:     []byte("*1\r\n,1.23\r\n"),
			expected: fmt.Sprintf("Array[Double{%s}]", strconv.FormatFloat(1.23, 'f', -1, 64)),
			isError:  false,
		},
		{
			name:     "array with array",
			a:        &Array{},
			line:     []byte("*1\r\n*1\r\n:123\r\n"),
			expected: "Array[Array[Integer{123}]]",
			isError:  false,
		},
		{
			name:     "array with multiple types",
			a:        &Array{},
			line:     []byte("*3\r\n:123\r\n+hello\r\n$5\r\nworld\r\n"),
			expected: "Array[Integer{123}, SimpleString{hello}, BulkString{world}]",
			isError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(string(tt.line)))
			err := tt.a.Read(reader)
			if tt.isError {
				if err == nil {
					t.Errorf("case %s: expected error but got nil", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("case %s: expected no error but got %v", tt.name, err)
				}
				if actual := tt.a.String(); actual != tt.expected {
					t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
				}
			}
		})
	}
}

func TestArray_Write(t *testing.T) {
	tests := []struct {
		name     string
		a        *Array
		expected string
	}{
		{
			name:     "empty array",
			a:        &Array{},
			expected: "*0\r\n",
		},
		{
			name:     "array with simplestring",
			a:        &Array{elements: []RedisObject{&SimpleString{value: "hello"}}},
			expected: "*1\r\n+hello\r\n",
		},
		{
			name:     "array with bulkstring",
			a:        &Array{elements: []RedisObject{&BulkString{value: []byte("hello")}}},
			expected: "*1\r\n$5\r\nhello\r\n",
		},
		{
			name:     "array with integer",
			a:        &Array{elements: []RedisObject{&Integer{value: 123}}},
			expected: "*1\r\n:123\r\n",
		},
		{
			name:     "array with null",
			a:        &Array{elements: []RedisObject{&Null{}}},
			expected: "*1\r\n_\r\n",
		},
		{
			name:     "array with boolean",
			a:        &Array{elements: []RedisObject{&Boolean{value: true}, &Boolean{value: false}}},
			expected: "*2\r\n#t\r\n#f\r\n",
		},
		{
			name:     "array with double",
			a:        &Array{elements: []RedisObject{&Double{value: 1.23}}},
			expected: fmt.Sprintf("*1\r\n,%s\r\n", strconv.FormatFloat(1.23, 'f', -1, 64)),
		},
		{
			name:     "array with array",
			a:        &Array{elements: []RedisObject{&Array{elements: []RedisObject{&Integer{value: 123}}}}},
			expected: "*1\r\n*1\r\n:123\r\n",
		},
		{
			name:     "array with multiple types",
			a:        &Array{elements: []RedisObject{&Integer{value: 123}, &SimpleString{value: "hello"}, &BulkString{value: []byte("world")}}},
			expected: "*3\r\n:123\r\n+hello\r\n$5\r\nworld\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := strings.Builder{}
			if err := tt.a.Write(&builder); err != nil {
				t.Errorf("case %s: unexpected error: %v", tt.name, err)
			}
			if actual := builder.String(); actual != tt.expected {
				t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestMap_Read(t *testing.T) {
	tests := []struct {
		name     string
		m        *Map
		line     []byte
		expected string
		isError  bool
	}{
		{
			name:     "empty map",
			m:        &Map{},
			line:     []byte("%0\r\n"),
			expected: "Map{}",
			isError:  false,
		},
		{
			name:    "malformed map size",
			m:       &Map{},
			line:    []byte("%a\r\n"),
			isError: true,
		},
		{
			name:    "map size insufficient",
			m:       &Map{},
			line:    []byte("%1\r\n"),
			isError: true,
		},
		{
			name:    "duplicated key map",
			m:       &Map{},
			line:    []byte("%2\r\n+hello\r\n+world\r\n+hello\r\n_\r\n"),
			isError: true,
		},
		{
			name:     "map with simple string",
			m:        &Map{},
			line:     []byte("%1\r\n+hello\r\n:123\r\n"),
			expected: "Map{hello: Integer{123}}",
			isError:  false,
		},
		{
			name:     "map with bulk string",
			m:        &Map{},
			line:     []byte("%1\r\n$5\r\nhello\r\n:123\r\n"),
			expected: "Map{hello: Integer{123}}",
			isError:  false,
		},
		{
			name:     "map with integer",
			m:        &Map{},
			line:     []byte("%1\r\n+123\r\n:123\r\n"),
			expected: "Map{123: Integer{123}}",
			isError:  false,
		},
		{
			name:     "map with null",
			m:        &Map{},
			line:     []byte("%1\r\n+nil\r\n_\r\n"),
			expected: "Map{nil: Null{}}",
			isError:  false,
		},
		{
			name:     "map with boolean",
			m:        &Map{},
			line:     []byte("%1\r\n+true\r\n#t\r\n"),
			expected: "Map{true: Boolean{true}}",
			isError:  false,
		},
		{
			name:     "map with double",
			m:        &Map{},
			line:     []byte("%1\r\n+1.23\r\n,1.23\r\n"),
			expected: fmt.Sprintf("Map{1.23: Double{%s}}", strconv.FormatFloat(1.23, 'f', -1, 64)),
			isError:  false,
		},
		{
			name:     "map with array",
			m:        &Map{},
			line:     []byte("%1\r\n+123\r\n*1\r\n:123\r\n"),
			expected: "Map{123: Array[Integer{123}]}",
			isError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(string(tt.line)))
			err := tt.m.Read(reader)
			if tt.isError {
				if err == nil {
					t.Errorf("case %s: expected error but got nil", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("case %s: expected no error but got %v", tt.name, err)
				}
				if actual := tt.m.String(); actual != tt.expected {
					t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
				}
			}
		})
	}
}

func TestMap_Write(t *testing.T) {
	tests := []struct {
		name     string
		m        *Map
		expected string
	}{
		{
			name:     "empty map",
			m:        &Map{},
			expected: "%0\r\n",
		},
		{
			name:     "map with simple string",
			m:        &Map{elements: map[string]RedisObject{"hello": &SimpleString{value: "hello"}}},
			expected: "%1\r\n$5\r\nhello\r\n+hello\r\n",
		},
		{
			name:     "map with bulk string",
			m:        &Map{elements: map[string]RedisObject{"hello": &BulkString{value: []byte("hello")}}},
			expected: "%1\r\n$5\r\nhello\r\n$5\r\nhello\r\n",
		},
		{
			name:     "map with integer",
			m:        &Map{elements: map[string]RedisObject{"123": &Integer{value: 123}}},
			expected: "%1\r\n$3\r\n123\r\n:123\r\n",
		},
		{
			name:     "map with null",
			m:        &Map{elements: map[string]RedisObject{"nil": &Null{}}},
			expected: "%1\r\n$3\r\nnil\r\n_\r\n",
		},
		{
			name:     "map with boolean",
			m:        &Map{elements: map[string]RedisObject{"true": &Boolean{value: true}}},
			expected: "%1\r\n$4\r\ntrue\r\n#t\r\n",
		},
		{
			name:     "map with double",
			m:        &Map{elements: map[string]RedisObject{"1.23": &Double{value: 1.23}}},
			expected: fmt.Sprintf("%%1\r\n$4\r\n1.23\r\n,%s\r\n", strconv.FormatFloat(1.23, 'f', -1, 64)),
		},
		{
			name:     "map with array",
			m:        &Map{elements: map[string]RedisObject{"123": &Array{elements: []RedisObject{&Integer{value: 123}}}}},
			expected: "%1\r\n$3\r\n123\r\n*1\r\n:123\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := strings.Builder{}
			if err := tt.m.Write(&builder); err != nil {
				t.Errorf("case %s: unexpected error: %v", tt.name, err)
			}
			if actual := builder.String(); actual != tt.expected {
				t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestSet_Read(t *testing.T) {
	tests := []struct {
		name     string
		s        *Set
		line     []byte
		expected string
		isError  bool
	}{
		{
			name:     "empty set",
			s:        &Set{},
			line:     []byte("~0\r\n"),
			expected: "Set{}",
			isError:  false,
		},
		{
			name:    "malformed set size",
			s:       &Set{},
			line:    []byte("~a\r\n"),
			isError: true,
		},
		{
			name:    "set size insufficient",
			s:       &Set{},
			line:    []byte("~1\r\n"),
			isError: true,
		},
		{
			name:    "duplicated value set",
			s:       &Set{},
			line:    []byte("~2\r\n+hello\r\n+hello\r\n"),
			isError: true,
		},
		{
			name:     "set with simple string",
			s:        &Set{},
			line:     []byte("~1\r\n+hello\r\n"),
			expected: "Set{SimpleString{hello}}",
			isError:  false,
		},
		{
			name:     "set with bulk string",
			s:        &Set{},
			line:     []byte("~1\r\n$5\r\nhello\r\n"),
			expected: "Set{BulkString{hello}}",
			isError:  false,
		},
		{
			name:     "set with integer",
			s:        &Set{},
			line:     []byte("~1\r\n:123\r\n"),
			expected: "Set{Integer{123}}",
			isError:  false,
		},
		{
			name:     "set with null",
			s:        &Set{},
			line:     []byte("~1\r\n_\r\n"),
			expected: "Set{Null{}}",
			isError:  false,
		},
		{
			name:     "set with boolean",
			s:        &Set{},
			line:     []byte("~1\r\n#t\r\n"),
			expected: "Set{Boolean{true}}",
			isError:  false,
		},
		{
			name:     "set with double",
			s:        &Set{},
			line:     []byte("~1\r\n,1.23\r\n"),
			expected: fmt.Sprintf("Set{Double{%s}}", strconv.FormatFloat(1.23, 'f', -1, 64)),
			isError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(string(tt.line)))
			err := tt.s.Read(reader)
			if tt.isError {
				if err == nil {
					t.Errorf("case %s: expected error but got nil", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("case %s: expected no error but got %v", tt.name, err)
				}
				if actual := tt.s.String(); actual != tt.expected {
					t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
				}
			}
		})
	}
}
