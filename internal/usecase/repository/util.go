package repository

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/entity"
	"github.com/jackc/pgx/v5/pgconn"
)

func selectRandomReviewers(members []entity.TeamMember, authorID string, count int) []string {
	var candidates []string
	for _, member := range members {
		if member.IsActive && member.ID != authorID {
			candidates = append(candidates, member.ID)
		}
	}

	if len(candidates) <= count {
		return candidates
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	return candidates[:count]
}

func findReplacementReviewer(
	members []entity.TeamMember,
	authorID, userID, currentReviewer1, currentReviewer2 string,
) string {
	var candidates []string
	for _, member := range members {
		if member.IsActive &&
			member.ID != authorID &&
			member.ID != userID &&
			member.ID != currentReviewer1 &&
			member.ID != currentReviewer2 {
			candidates = append(candidates, member.ID)
		}
	}

	if len(candidates) == 0 {
		return ""
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	return candidates[0]
}

func checkUniqueKeyViolation(err, errExist error, id string) error {
	const ErrUniqueKeyViolation = "23505"
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == ErrUniqueKeyViolation {
		return fmt.Errorf("ID %s exists: %w",
			id, errExist)
	}
	return err
}
