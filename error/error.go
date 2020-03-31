package error

import "errors"

var (
	AliasNotFoundError = errors.New("could not find alias")
)
