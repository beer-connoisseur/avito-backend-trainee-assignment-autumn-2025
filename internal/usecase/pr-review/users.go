package pr_review

import (
	"context"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/entity"
)

func (p *prReviewImpl) SetUserActive(ctx context.Context, userID string, isActive bool) (*entity.User, error) {
	user, err := p.usersRepository.SetUserActive(ctx, userID, isActive)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *prReviewImpl) GetUserReviews(ctx context.Context, userID string) (string, []*entity.PullRequestShort, error) {
	ID, pr, err := p.usersRepository.GetUserReviews(ctx, userID)
	if err != nil {
		return "", nil, err
	}

	return ID, pr, nil
}
