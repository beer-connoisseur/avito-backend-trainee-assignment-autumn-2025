package controller

import (
	"context"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/generated/api/pr-review"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *implementation) GetUserReviews(ctx context.Context, req *pr_review.GetUserReviewsRequest) (*pr_review.GetUserReviewsResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, prs, err := i.usersUseCase.GetUserReviews(ctx, req.GetUserId())
	if err != nil {
		return nil, i.convertErr(err)
	}

	return &pr_review.GetUserReviewsResponse{
		UserId:       id,
		PullRequests: i.convertPullRequestsShortToProto(prs),
	}, nil
}
