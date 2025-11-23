package pr_review

import (
	"context"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/entity"
	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/usecase/repository"
	"go.uber.org/zap"
)

type (
	UsersUseCase interface {
		SetUserActive(_ context.Context, userID string, isActive bool) (*entity.User, error)
		GetUserReviews(_ context.Context, userID string) (string, []*entity.PullRequestShort, error)
	}

	TeamsUseCase interface {
		AddTeam(_ context.Context, name string, members []entity.TeamMember) (*entity.Team, error)
		GetTeam(_ context.Context, name string) (*entity.Team, error)
	}

	PullRequestsUseCase interface {
		CreatePullRequest(_ context.Context, prID, name, authorID string) (*entity.PullRequest, error)
		MergePullRequest(_ context.Context, prID string) (*entity.PullRequest, error)
		ReassignReviewer(_ context.Context, prID, userID string) (*entity.PullRequest, string, error)
	}
)

var _ UsersUseCase = (*prReviewImpl)(nil)
var _ TeamsUseCase = (*prReviewImpl)(nil)
var _ PullRequestsUseCase = (*prReviewImpl)(nil)

type prReviewImpl struct {
	logger                 *zap.Logger
	usersRepository        repository.UsersRepository
	teamsRepository        repository.TeamsRepository
	pullRequestsRepository repository.PullRequestsRepository
}

func New(
	logger *zap.Logger,
	usersRepository repository.UsersRepository,
	teamsRepository repository.TeamsRepository,
	pullRequestsRepository repository.PullRequestsRepository,
) *prReviewImpl {
	return &prReviewImpl{
		logger:                 logger,
		usersRepository:        usersRepository,
		teamsRepository:        teamsRepository,
		pullRequestsRepository: pullRequestsRepository,
	}
}
