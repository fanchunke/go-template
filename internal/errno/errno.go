package errno

type CustomError interface {
	error
	WithError(err error) CustomError
}

type customError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (e customError) Error() string {
	return e.Message
}

func (e customError) WithError(err error) CustomError {
	e.Detail = err.Error()
	return e
}

func New(code int, message string) CustomError {
	return &customError{
		Code:    code,
		Message: message,
	}
}
