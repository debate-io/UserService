package userepo

import "github.com/ztrue/tracerr"

var (
	ErrUserAlreadyExist = tracerr.New("user already register")
	ErrUserNotFound     = tracerr.New("user not found")
	//role?
)
