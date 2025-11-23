package controller

import (
	"context"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/generated/api/pr-review"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *implementation) SetUserActive(ctx context.Context, req *pr_review.SetUserActiveRequest) (*pr_review.SetUserActiveResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := i.usersUseCase.SetUserActive(ctx, req.GetUserId(), req.GetIsActive())
	if err != nil {
		return nil, i.convertErr(err)
	}

	return &pr_review.SetUserActiveResponse{
		User: &pr_review.User{
			UserId:   user.ID,
			Username: user.Name,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	}, nil
}
