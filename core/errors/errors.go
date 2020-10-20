package errors

import "errors"

// ApplicationError ...
type ApplicationError struct {
	Err    error
	Msg    string
	Params map[string]interface{}
}

// NewApplicationError ...
func NewApplicationError(msg string, params map[string]interface{}) ApplicationError {
	return ApplicationError{
		Msg:    msg,
		Params: params,
		Err:    errors.New(msg),
	}
}

func (e ApplicationError) Error() string {
	return e.Err.Error()
}

// TBD: the names shold be revisited according to the standards
var (
	UnauthorizedErr = ApplicationError{Msg: "Unauthorized Access"}
	ApplicationErr  = ApplicationError{Msg: "Application error"}
)
