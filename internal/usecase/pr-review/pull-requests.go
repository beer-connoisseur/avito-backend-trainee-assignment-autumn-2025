package pr_review

import (
	"context"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/entity"
)

func (p *prReviewImpl) CreatePullRequest(ctx context.Context, prID, name, authorID string) (*entity.PullRequest, error) {
	pr, err := p.pullRequestsRepository.CreatePullRequest(ctx, prID, name, authorID)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (p *prReviewImpl) MergePullRequest(ctx context.Context, prID string) (*entity.PullRequest, error) {
	pr, err := p.pullRequestsRepository.MergePullRequest(ctx, prID)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (p *prReviewImpl) ReassignReviewer(ctx context.Context, prID, userID string) (*entity.PullRequest, string, error) {
	pr, newID, err := p.pullRequestsRepository.ReassignReviewer(ctx, prID, userID)
	if err != nil {
		return nil, "", err
	}

	return pr, newID, nil
}
