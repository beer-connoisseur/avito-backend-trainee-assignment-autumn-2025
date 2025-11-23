package controller

import (
	"errors"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/generated/api/pr-review"
	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *implementation) convertProtoMembersToEntity(protoMembers []*pr_review.TeamMember) []entity.TeamMember {
	members := make([]entity.TeamMember, len(protoMembers))
	for j, protoMember := range protoMembers {
		members[j] = entity.TeamMember{
			ID:       protoMember.GetUserId(),
			Name:     protoMember.GetUsername(),
			IsActive: protoMember.GetIsActive(),
		}
	}
	return members
}

func (i *implementation) convertMembersToProto(team *entity.Team) []*pr_review.TeamMember {
	protoMembers := make([]*pr_review.TeamMember, len(team.Members))
	for j, member := range team.Members {
		protoMembers[j] = &pr_review.TeamMember{
			UserId:   member.ID,
			Username: member.Name,
			IsActive: member.IsActive,
		}
	}

	return protoMembers
}

func (i *implementation) convertPrStatusToProto(status string) pr_review.PRStatus {
	if status == "MERGED" {
		return pr_review.PRStatus_MERGED
	}
	return pr_review.PRStatus_OPEN
}

func (i *implementation) convertPullRequestsShortToProto(prs []*entity.PullRequestShort) []*pr_review.PullRequestShort {
	protoPRs := make([]*pr_review.PullRequestShort, len(prs))
	for j, pr := range prs {
		protoPRs[j] = &pr_review.PullRequestShort{
			PullRequestId:   pr.ID,
			PullRequestName: pr.Name,
			AuthorId:        pr.AuthorID,
			Status:          i.convertPrStatusToProto(pr.Status),
		}
	}
	return protoPRs
}

func (i *implementation) convertErr(err error) error {
	switch {
	case errors.Is(err, entity.ErrPrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, entity.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, entity.ErrTeamNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, entity.ErrPrMerged):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, entity.ErrPrNoCandidate):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, entity.ErrUserNotAssigned):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, entity.ErrPrExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, entity.ErrTeamExists):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, entity.ErrTeamNoMembers):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
