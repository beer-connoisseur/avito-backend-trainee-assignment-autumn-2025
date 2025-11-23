package repository

import (
	"context"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/entity"
)

type (
	UsersRepository interface {
		SetUserActive(_ context.Context, userID string, isActive bool) (*entity.User, error)
		GetUserReviews(_ context.Context, userID string) (string, []*entity.PullRequestShort, error)
	}

	TeamsRepository interface {
		AddTeam(_ context.Context, name string, members []entity.TeamMember) (*entity.Team, error)
		GetTeam(_ context.Context, name string) (*entity.Team, error)
	}

	PullRequestsRepository interface {
		CreatePullRequest(_ context.Context, prID, name, authorID string) (*entity.PullRequest, error)
		MergePullRequest(_ context.Context, prID string) (*entity.PullRequest, error)
		ReassignReviewer(_ context.Context, prID, userID string) (*entity.PullRequest, string, error)
	}
)
