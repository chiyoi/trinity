package atmt

import "fmt"

type PostError interface {
	error
	StatusCode() StatusCode
}

type postError struct {
	code StatusCode
}

var _ PostError = (*postError)(nil)

func (err *postError) Error() string {
	return fmt.Sprintf("post error(%d %s)", err.code, err.code)
}

func (err *postError) StatusCode() StatusCode { return err.code }

type messageTypeError struct {
	typ, exp MessageType
}

var _ error = (*messageTypeError)(nil)

func (err *messageTypeError) Error() string {
	return fmt.Sprintf("unexpected message type(%s), expected(%s)", err.typ, err.exp)
}
