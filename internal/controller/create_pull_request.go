package controller

import (
	"context"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/generated/api/pr-review"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *implementation) CreatePullRequest(ctx context.Context, req *pr_review.CreatePullRequestRequest) (*pr_review.CreatePullRequestResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pr, err := i.pullRequestsUseCase.CreatePullRequest(ctx, req.GetPullRequestId(), req.GetPullRequestName(), req.GetAuthorId())
	if err != nil {
		return nil, i.convertErr(err)
	}

	response := &pr_review.CreatePullRequestResponse{
		Pr: &pr_review.PullRequest{
			PullRequestId:     pr.ID,
			PullRequestName:   pr.Name,
			AuthorId:          pr.AuthorID,
			Status:            i.convertPrStatusToProto(pr.Status),
			AssignedReviewers: []string{pr.AssignedReviewerFirst, pr.AssignedReviewerSecond},
			CreatedAt:         timestamppb.New(*pr.CreatedAt),
		},
	}

	if pr.MergedAt != nil {
		response.Pr.MergedAt = timestamppb.New(*pr.MergedAt)
	}

	return response, nil
}
