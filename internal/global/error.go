package global

import "errors"

var (
	ErrorForbidden      = errors.New("forbidden")
	ErrorNotFound       = errors.New("not found")
	ErrorInternalServer = errors.New("internal server error")
)
