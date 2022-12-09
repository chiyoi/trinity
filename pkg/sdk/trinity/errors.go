package trinity

import (
	"fmt"

	"github.com/chiyoi/trinity/pkg/atmt"
)

type RequestError interface {
	error
	Action() Action
	StatusCode() atmt.StatusCode
}

type requestError struct {
	act  Action
	code atmt.StatusCode
}

var _ RequestError = (*requestError)(nil)

func (err *requestError) Error() string {
	return fmt.Sprintf("post error(%d %s)", err.code, err.code.Text())
}

func (err *requestError) Action() Action              { return err.act }
func (err *requestError) StatusCode() atmt.StatusCode { return err.code }
