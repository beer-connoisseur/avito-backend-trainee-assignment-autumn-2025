package entity

import (
	"errors"
	"time"
)

type PullRequest struct {
	ID                     string
	Name                   string
	AuthorID               string
	Status                 string
	AssignedReviewerFirst  string
	AssignedReviewerSecond string
	CreatedAt              *time.Time
	MergedAt               *time.Time
}

type PullRequestShort struct {
	ID       string
	Name     string
	AuthorID string
	Status   string
}

const (
	OPEN   = "OPEN"
	MERGED = "MERGED"
)

var (
	ErrPrNoCandidate = errors.New("no active replacement candidate")
	ErrPrExists      = errors.New("pull request already exists")
	ErrPrMerged      = errors.New("pull request is merged")
	ErrPrNotFound    = errors.New("pull request was not found")
)
