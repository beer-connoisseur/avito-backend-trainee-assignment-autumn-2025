package entity

import "errors"

type User struct {
	ID       string
	Name     string
	TeamName string
	IsActive bool
}

type TeamMember struct {
	ID       string
	Name     string
	IsActive bool
}

var (
	ErrUserNotFound    = errors.New("user was not found")
	ErrUserNotAssigned = errors.New("reviewer not assigned")
)
