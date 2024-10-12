package repo

import "github.com/ztrue/tracerr"

var (
	ErrUserAlreadyExist     = tracerr.New("user already register")
	ErrUserNotFound         = tracerr.New("user not found")
	ErrRecoveryCodeNotFound = tracerr.New("user code not found")

	//role?
)
