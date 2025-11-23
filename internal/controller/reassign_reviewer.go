package controller

import (
	"context"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/generated/api/pr-review"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (i *implementation) ReassignReviewer(ctx context.Context, req *pr_review.ReassignReviewerRequest) (*pr_review.ReassignReviewerResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pr, newReviewer, err := i.pullRequestsUseCase.ReassignReviewer(ctx, req.GetPullRequestId(), req.GetOldReviewerId())
	if err != nil {
		return nil, i.convertErr(err)
	}

	response := &pr_review.ReassignReviewerResponse{
		Pr: &pr_review.PullRequest{
			PullRequestId:     pr.ID,
			PullRequestName:   pr.Name,
			AuthorId:          pr.AuthorID,
			Status:            i.convertPrStatusToProto(pr.Status),
			AssignedReviewers: []string{pr.AssignedReviewerFirst, pr.AssignedReviewerSecond},
			CreatedAt:         timestamppb.New(*pr.CreatedAt),
		},
		ReplacedBy: newReviewer,
	}

	if pr.MergedAt != nil {
		response.Pr.MergedAt = timestamppb.New(*pr.MergedAt)
	}

	return response, nil
}
