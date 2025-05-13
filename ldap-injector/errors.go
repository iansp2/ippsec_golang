package main

type PasswordError struct {
	Err  error
	Code int
}

func (e *PasswordError) Error() string {
	return e.Err.Error()
}

func NewPasswordErrorWithCode(err error, code int) error {
	return &PasswordError{Err: err, Code: code}
}

func NewPasswordError(err error) error {
	return &PasswordError{Err: err}
}
