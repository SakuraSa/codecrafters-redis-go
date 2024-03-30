package redis

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestSimpleString_Read(t *testing.T) {
	tests := []struct {
		name     string
		s        *SimpleString
		line     []byte
		expected string
		isError  bool
	}{
		{
			name:     "simple string",
			s:        &SimpleString{},
			line:     []byte("+OK\r\n"),
			expected: "SimpleString{OK}",
		},
		{
			name:     "unexpected leading",
			s:        &SimpleString{},
			line:     []byte("-OK\r\n"),
			isError:  true,
			expected: "",
		},
		{
			name:     "empty string",
			s:        &SimpleString{},
			line:     []byte("+\r\n"),
			isError:  false,
			expected: "SimpleString{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(string(tt.line)))
			err := tt.s.Read(reader)
			if tt.isError && err == nil {
				t.Errorf("case %s: expected error, but got nil", tt.name)
			}
			if !tt.isError {
				if err != nil {
					t.Errorf("case %s: unexpected error: %v", tt.name, err)
				} else if actual := tt.expected; actual != tt.s.String() {
					t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
				}
			}
		})
	}
}

func TestSimpleString_Write(t *testing.T) {
	tests := []struct {
		name     string
		s        *SimpleString
		expected string
	}{
		{
			name:     "simple string",
			s:        &SimpleString{value: "OK"},
			expected: "+OK\r\n",
		},
		{
			name:     "empty string",
			s:        &SimpleString{},
			expected: "+\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &strings.Builder{}
			err := tt.s.Write(writer)
			if err != nil {
				t.Errorf("case %s: unexpected error: %v", tt.name, err)
			}
			if actual := writer.String(); actual != tt.expected {
				t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestSimpleError_Read(t *testing.T) {
	tests := []struct {
		name     string
		s        *SimpleError
		line     []byte
		expected string
		isError  bool
	}{
		{
			name:     "simple error",
			s:        &SimpleError{},
			line:     []byte("-ERR unknown command\r\n"),
			expected: "SimpleError{ERR unknown command}",
		},
		{
			name:     "unexpected leading",
			s:        &SimpleError{},
			line:     []byte("+ERR unknown command\r\n"),
			isError:  true,
			expected: "",
		},
		{
			name:     "empty string",
			s:        &SimpleError{},
			line:     []byte("-\r\n"),
			isError:  false,
			expected: "SimpleError{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(string(tt.line)))
			err := tt.s.Read(reader)
			if tt.isError && err == nil {
				t.Errorf("case %s: expected error, but got nil", tt.name)
			}
			if !tt.isError {
				if err != nil {
					t.Errorf("case %s: unexpected error: %v", tt.name, err)
				} else if actual := tt.s.String(); actual != tt.expected {
					t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
				}
			}
		})
	}
}

func TestSimpleError_Write(t *testing.T) {
	tests := []struct {
		name     string
		s        *SimpleError
		expected string
	}{
		{
			name:     "simple error",
			s:        &SimpleError{value: "ERR unknown command"},
			expected: "-ERR unknown command\r\n",
		},
		{
			name:     "empty string",
			s:        &SimpleError{},
			expected: "-\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &strings.Builder{}
			err := tt.s.Write(writer)
			if err != nil {
				t.Errorf("case %s: unexpected error: %v", tt.name, err)
			}
			if actual := writer.String(); actual != tt.expected {
				t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestInteger_Read(t *testing.T) {
	tests := []struct {
		name     string
		i        *Integer
		line     []byte
		expected string
		isError  bool
	}{
		{
			name:     "integer",
			i:        &Integer{},
			line:     []byte(":1000\r\n"),
			expected: "Integer{1000}",
		},
		{
			name:     "unexpected leading",
			i:        &Integer{},
			line:     []byte("$1000\r\n"),
			isError:  true,
			expected: "",
		},
		{
			name:     "minus integer",
			i:        &Integer{},
			line:     []byte(":-1000\r\n"),
			isError:  false,
			expected: "Integer{-1000}",
		},
		{
			name:     "malformed integer",
			i:        &Integer{},
			line:     []byte(":not_int\r\n"),
			isError:  true,
			expected: "",
		},
		{
			name:     "empty string",
			i:        &Integer{},
			line:     []byte(":\r\n"),
			isError:  true,
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(string(tt.line)))
			err := tt.i.Read(reader)
			if tt.isError && err == nil {
				t.Errorf("case %s: expected error, but got nil", tt.name)
			}
			if !tt.isError {
				if err != nil {
					t.Errorf("case %s: unexpected error: %v", tt.name, err)
				} else if actual := tt.i.String(); actual != tt.expected {
					t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
				}
			}

		})
	}
}

func TestInteger_Write(t *testing.T) {
	tests := []struct {
		name     string
		i        *Integer
		expected string
	}{
		{
			name:     "integer",
			i:        &Integer{value: 1000},
			expected: ":1000\r\n",
		},
		{
			name:     "minus integer",
			i:        &Integer{value: -1000},
			expected: ":-1000\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &strings.Builder{}
			err := tt.i.Write(writer)
			if err != nil {
				t.Errorf("case %s: unexpected error: %v", tt.name, err)
			}
			if actual := writer.String(); actual != tt.expected {
				t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestBulkString_Read(t *testing.T) {
	tests := []struct {
		name     string
		b        *BulkString
		line     []byte
		expected string
		isError  bool
	}{
		{
			name:     "bulk string",
			b:        &BulkString{},
			line:     []byte("$6\r\nfoobar\r\n"),
			expected: "BulkString{foobar}",
		},
		{
			name:     "unexpected leading",
			b:        &BulkString{},
			line:     []byte("#6\r\nfoobar\r\n"),
			isError:  true,
			expected: "",
		},
		{
			name:     "empty string",
			b:        &BulkString{},
			line:     []byte("$0\r\n\r\n"),
			isError:  false,
			expected: "BulkString{}",
		},
		{
			name:     "null string",
			b:        &BulkString{},
			line:     []byte("$-1\r\n"),
			isError:  false,
			expected: "BulkString{}",
		},
		{
			name:     "malformed bulk string",
			b:        &BulkString{},
			line:     []byte("$6\r\nfoo\r\n"),
			isError:  true,
			expected: "",
		},
		{
			name:     "string with endline",
			b:        &BulkString{},
			line:     []byte("$9\r\nendline\r\n\r\n"),
			isError:  false,
			expected: "BulkString{endline\r\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(string(tt.line)))
			err := tt.b.Read(reader)
			if tt.isError && err == nil {
				t.Errorf("case %s: expected error, but got nil", tt.name)
			}
			if !tt.isError {
				if err != nil {
					t.Errorf("case %s: unexpected error: %v", tt.name, err)
				} else if actual := tt.b.String(); actual != tt.expected {
					t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
				}
			}

		})
	}
}

func TestBulkString_Write(t *testing.T) {
	tests := []struct {
		name     string
		b        *BulkString
		expected string
	}{
		{
			name:     "bulk string",
			b:        &BulkString{value: []byte("foobar")},
			expected: "$6\r\nfoobar\r\n",
		},
		{
			name:     "empty string",
			b:        &BulkString{value: []byte{}},
			expected: "$0\r\n\r\n",
		},
		{
			name:     "null string",
			b:        &BulkString{value: nil},
			expected: "$-1\r\n",
		},
		{
			name:     "string with endline",
			b:        &BulkString{value: []byte("endline\r\n")},
			expected: "$9\r\nendline\r\n\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &strings.Builder{}
			err := tt.b.Write(writer)
			if err != nil {
				t.Errorf("case %s: unexpected error: %v", tt.name, err)
			}
			if actual := writer.String(); actual != tt.expected {
				t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestNull_Read(t *testing.T) {
	tests := []struct {
		name    string
		n       *Null
		line    []byte
		isError bool
	}{
		{
			name:    "unexpected leading",
			n:       &Null{},
			line:    []byte("#-1\r\n"),
			isError: true,
		},
		{
			name:    "malformed null",
			n:       &Null{},
			line:    []byte("_-2\r\n"),
			isError: true,
		},
		{
			name:    "null",
			n:       &Null{},
			line:    []byte("_\r\n"),
			isError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(string(tt.line)))
			err := tt.n.Read(reader)
			if tt.isError && err == nil {
				t.Errorf("case %s: expected error, but got nil", tt.name)
			}
			if !tt.isError && err != nil {
				t.Errorf("case %s: unexpected error: %v", tt.name, err)
			}
		})
	}
}

func TestNull_Write(t *testing.T) {
	tests := []struct {
		name     string
		n        *Null
		expected string
	}{
		{
			name:     "null",
			n:        &Null{},
			expected: "_\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &strings.Builder{}
			err := tt.n.Write(writer)
			if err != nil {
				t.Errorf("case %s: unexpected error: %v", tt.name, err)
			}
			if actual := writer.String(); actual != tt.expected {
				t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestBoolean_Read(t *testing.T) {
	tests := []struct {
		name     string
		b        *Boolean
		line     []byte
		isError  bool
		expected string
	}{
		{
			name:     "unexpected leading",
			b:        &Boolean{},
			line:     []byte("$true\r\n"),
			isError:  true,
			expected: "",
		},
		{
			name:     "malformed boolean",
			b:        &Boolean{},
			line:     []byte("#true\r\n"),
			isError:  true,
			expected: "",
		},
		{
			name:     "boolean true",
			b:        &Boolean{},
			line:     []byte("#t\r\n"),
			isError:  false,
			expected: "Boolean{true}",
		},
		{
			name:     "boolean false",
			b:        &Boolean{},
			line:     []byte("#f\r\n"),
			isError:  false,
			expected: "Boolean{false}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(string(tt.line)))
			err := tt.b.Read(reader)
			if tt.isError && err == nil {
				t.Errorf("case %s: expected error, but got nil", tt.name)
			}
			if !tt.isError {
				if err != nil {
					t.Errorf("case %s: unexpected error: %v", tt.name, err)
				} else if actual := tt.b.String(); actual != tt.expected {
					t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
				}
			}

		})
	}
}

func TestBoolean_Write(t *testing.T) {
	tests := []struct {
		name     string
		b        *Boolean
		expected string
	}{
		{
			name:     "boolean true",
			b:        &Boolean{value: true},
			expected: "#t\r\n",
		},
		{
			name:     "boolean false",
			b:        &Boolean{value: false},
			expected: "#f\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &strings.Builder{}
			err := tt.b.Write(writer)
			if err != nil {
				t.Errorf("case %s: unexpected error: %v", tt.name, err)
			}
			if actual := writer.String(); actual != tt.expected {
				t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestDouble_Read(t *testing.T) {
	tests := []struct {
		name     string
		d        *Double
		line     []byte
		isError  bool
		expected string
	}{
		{
			name:     "unexpected leading",
			d:        &Double{},
			line:     []byte("$3.14\r\n"),
			isError:  true,
			expected: "",
		},
		{
			name:     "malformed double",
			d:        &Double{},
			line:     []byte(",not_double\r\n"),
			isError:  true,
			expected: "",
		},
		{
			name:     "double",
			d:        &Double{},
			line:     []byte(",3.14\r\n"),
			isError:  false,
			expected: fmt.Sprintf("Double{%s}", strconv.FormatFloat(3.14, 'f', -1, 64)),
		},
		{
			name:     "minus double",
			d:        &Double{},
			line:     []byte(",-3.14\r\n"),
			isError:  false,
			expected: fmt.Sprintf("Double{%s}", strconv.FormatFloat(-3.14, 'f', -1, 64)),
		},
		{
			name:     "exp double",
			d:        &Double{},
			line:     []byte(",-3.14e-3\r\n"),
			isError:  false,
			expected: fmt.Sprintf("Double{%s}", strconv.FormatFloat(-3.14e-3, 'f', -1, 64)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(string(tt.line)))
			err := tt.d.Read(reader)
			if tt.isError && err == nil {
				t.Errorf("case %s: expected error, but got nil", tt.name)
			}
			if !tt.isError {
				if err != nil {
					t.Errorf("case %s: unexpected error: %v", tt.name, err)
				} else if actual := tt.d.String(); actual != tt.expected {
					t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
				}
			}

		})
	}
}

func TestDouble_Write(t *testing.T) {
	tests := []struct {
		name     string
		d        *Double
		expected string
	}{
		{
			name:     "double",
			d:        &Double{value: 3.14},
			expected: ",3.14\r\n",
		},
		{
			name:     "minus double",
			d:        &Double{value: -3.14},
			expected: ",-3.14\r\n",
		},
		{
			name:     "exp double",
			d:        &Double{value: -3.14e-3},
			expected: ",-0.00314\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &strings.Builder{}
			err := tt.d.Write(writer)
			if err != nil {
				t.Errorf("case %s: unexpected error: %v", tt.name, err)
			} else if actual := writer.String(); actual != tt.expected {
				t.Errorf("case %s: expected=%s, actual=%s", tt.name, tt.expected, actual)
			}
		})
	}
}
