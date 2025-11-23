package controller

import (
	"context"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/generated/api/pr-review"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *implementation) AddTeam(ctx context.Context, req *pr_review.AddTeamRequest) (*pr_review.AddTeamResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	team, err := i.teamsUseCase.AddTeam(ctx, req.GetTeamName(), i.convertProtoMembersToEntity(req.GetMembers()))
	if err != nil {
		return nil, i.convertErr(err)
	}

	return &pr_review.AddTeamResponse{
		Team: &pr_review.Team{
			TeamName: team.Name,
			Members:  i.convertMembersToProto(team),
		},
	}, nil
}
