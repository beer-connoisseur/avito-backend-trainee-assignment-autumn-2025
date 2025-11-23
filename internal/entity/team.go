package entity

import "errors"

type Team struct {
	Name    string
	Members []TeamMember
}

var (
	ErrTeamNoMembers = errors.New("no team members found")
	ErrTeamExists    = errors.New("team already exists")
	ErrTeamNotFound  = errors.New("team was not found")
)
