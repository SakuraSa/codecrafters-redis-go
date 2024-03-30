package redis

import "fmt"

var (
	ErrUnexpectedLeading = &UnexpectedLeadingError{}
	ErrUnexpectedTailing = &UnexpectedTailingError{}
	ErrSyntaxError       = &SyntaxError{}
)

type UnexpectedLeadingError struct {
	Expected byte
	Actual   byte
}

func (e *UnexpectedLeadingError) Error() string {
	return fmt.Sprintf("unexpected leading byte expected=%d\"%c\" actual=%d\"%c\"",
		e.Expected, e.Expected, e.Actual, e.Actual)
}

type UnexpectedTailingError struct {
	Expected string
	Actual   string
}

func (e *UnexpectedTailingError) Error() string {
	return fmt.Sprintf("unexpected tailing byte expected=%v actual=%v",
		e.Expected, e.Actual)
}

type SyntaxError struct {
	Msg string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("syntax error: %s", e.Msg)
}
