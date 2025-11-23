package controller

import (
	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/usecase/pr-review"
	"go.uber.org/zap"
)

type implementation struct {
	logger              *zap.Logger
	usersUseCase        pr_review.UsersUseCase
	teamsUseCase        pr_review.TeamsUseCase
	pullRequestsUseCase pr_review.PullRequestsUseCase
}

func New(
	logger *zap.Logger,
	usersUseCase pr_review.UsersUseCase,
	teamsUseCase pr_review.TeamsUseCase,
	pullRequestsUseCase pr_review.PullRequestsUseCase,
) *implementation {
	return &implementation{
		logger:              logger,
		usersUseCase:        usersUseCase,
		teamsUseCase:        teamsUseCase,
		pullRequestsUseCase: pullRequestsUseCase,
	}
}
