package pr_review

import (
	"context"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/entity"
)

func (p *prReviewImpl) AddTeam(ctx context.Context, name string, members []entity.TeamMember) (*entity.Team, error) {
	team, err := p.teamsRepository.AddTeam(ctx, name, members)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (p *prReviewImpl) GetTeam(ctx context.Context, name string) (*entity.Team, error) {
	team, err := p.teamsRepository.GetTeam(ctx, name)
	if err != nil {
		return nil, err
	}

	return team, nil
}
