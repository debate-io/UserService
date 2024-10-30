package repo

import "github.com/ztrue/tracerr"

var (
	ErrAlreadyExist = tracerr.New("already exists")
	ErrNotFound     = tracerr.New("not found")
)
