package controller

import (
	"context"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/generated/api/pr-review"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *implementation) GetTeam(ctx context.Context, req *pr_review.GetTeamRequest) (*pr_review.GetTeamResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	team, err := i.teamsUseCase.GetTeam(ctx, req.GetTeamName())
	if err != nil {
		return nil, i.convertErr(err)
	}

	return &pr_review.GetTeamResponse{
		Team: &pr_review.Team{
			TeamName: team.Name,
			Members:  i.convertMembersToProto(team),
		},
	}, nil
}
