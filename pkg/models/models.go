package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("Models: no matching record found")
	ErrInvalidCredentials = errors.New("Models: Invalid credentials")
	ErrDuplicateEmail     = errors.New("Models: Duplicate Email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}
