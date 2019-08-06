package database

import "fmt"

type UserAlreadyExistsError struct {
	err      string
	username string
}

func (e *UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user %s already exists: %s", e.username, e.err)
}
