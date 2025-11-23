package repository

import (
	"context"
	"errors"
	"time"

	"github.com/beer-connoisseur/avito-backend-trainee-assignment-autumn-2025/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const reviewersCount = 2

type postgresRepository struct {
	logger *zap.Logger
	db     *pgxpool.Pool
}

func NewPostgresRepository(logger *zap.Logger, db *pgxpool.Pool) *postgresRepository {
	return &postgresRepository{
		logger: logger,
		db:     db,
	}
}

func (p *postgresRepository) SetUserActive(ctx context.Context, userID string, isActive bool) (*entity.User, error) {
	const query = `
UPDATE users
SET is_active = $1
WHERE id = $2
RETURNING id, name, (SELECT name FROM teams WHERE user_id = $2 LIMIT 1) as team_name, is_active
`
	var user entity.User
	err := p.db.QueryRow(ctx, query, isActive, userID).Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &entity.User{}, entity.ErrUserNotFound
		}
		return &entity.User{}, err
	}

	return &user, nil
}

func (p *postgresRepository) GetUserReviews(ctx context.Context, userID string) (string, []*entity.PullRequestShort, error) {
	const query = `
SELECT id, name, author_id, status
FROM pull_requests 
WHERE (assigned_reviewer_first = $1 OR assigned_reviewer_second = $1)
AND status = 'OPEN'
`
	rows, err := p.db.Query(ctx, query, userID)
	if err != nil {
		return userID, nil, err
	}
	defer rows.Close()

	var prs []*entity.PullRequestShort
	for rows.Next() {
		var pr entity.PullRequestShort
		if err = rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status); err != nil {
			return userID, nil, err
		}
		prs = append(prs, &pr)
	}

	if err = rows.Err(); err != nil {
		return userID, nil, err
	}

	return userID, prs, nil
}

func (p *postgresRepository) AddTeam(ctx context.Context, name string, members []entity.TeamMember) (*entity.Team, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return &entity.Team{}, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
		if err != nil {
			p.logger.Debug("Error: rollback")
		}
	}(tx, ctx)

	if len(members) == 0 {
		return &entity.Team{}, entity.ErrTeamNoMembers
	}

	exists, err := p.teamExists(ctx, tx, name)
	if err != nil {
		return &entity.Team{}, err
	}
	if exists {
		return &entity.Team{}, entity.ErrTeamExists
	}

	const queryUser = `
INSERT INTO users (id, name)
VALUES ($1, $2)
`
	const queryTeam = `
INSERT INTO teams (name, user_id)
VALUES ($1, $2)
`
	for _, member := range members {
		_, err = tx.Exec(ctx, queryUser, member.ID, member.Name)
		if err != nil {
			return &entity.Team{}, err
		}

		_, err = tx.Exec(ctx, queryTeam, name, member.ID)
		if err != nil {
			return &entity.Team{}, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return &entity.Team{}, err
	}

	team := &entity.Team{
		Name:    name,
		Members: members,
	}

	return team, nil
}

func (p *postgresRepository) GetTeam(ctx context.Context, name string) (*entity.Team, error) {
	const query = `
SELECT u.id, u.name, u.is_active 
FROM users u
JOIN teams t ON u.id = t.user_id
WHERE t.name = $1
`
	rows, err := p.db.Query(ctx, query, name)
	if err != nil {
		return &entity.Team{}, err
	}
	defer rows.Close()

	team := &entity.Team{
		Name: name,
	}

	for rows.Next() {
		var member entity.TeamMember
		err = rows.Scan(&member.ID, &member.Name, &member.IsActive)
		if err != nil {
			return &entity.Team{}, err
		}
		team.Members = append(team.Members, member)
	}

	if len(team.Members) == 0 {
		return &entity.Team{}, entity.ErrTeamNotFound
	}

	if err = rows.Err(); err != nil {
		return &entity.Team{}, err
	}

	return team, nil
}

func (p *postgresRepository) CreatePullRequest(ctx context.Context, prID, name, authorID string) (*entity.PullRequest, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return &entity.PullRequest{}, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
		if err != nil {
			p.logger.Debug("Error: rollback")
		}
	}(tx, ctx)

	team, err := p.getTeamByUserID(ctx, tx, authorID)
	if err != nil {
		return &entity.PullRequest{}, err
	}

	reviewers := selectRandomReviewers(team.Members, authorID, reviewersCount)

	reviewer1, reviewer2 := "", ""
	if len(reviewers) > 0 {
		reviewer1 = reviewers[0]
	}
	if len(reviewers) > 1 {
		reviewer2 = reviewers[1]
	}

	const query = `
INSERT INTO pull_requests (id, name, author_id, assigned_reviewer_first, assigned_reviewer_second)
VALUES ($1, $2, $3, $4, $5)
RETURNING created_at
`
	var createdAt *time.Time
	err = tx.QueryRow(ctx, query, prID, name, authorID, reviewer1, reviewer2).Scan(&createdAt)
	if err != nil {
		return &entity.PullRequest{}, checkUniqueKeyViolation(err, entity.ErrPrExists, prID)
	}

	if err = tx.Commit(ctx); err != nil {
		return &entity.PullRequest{}, err
	}

	pr := &entity.PullRequest{
		ID:                     prID,
		Name:                   name,
		AuthorID:               authorID,
		Status:                 entity.OPEN,
		AssignedReviewerFirst:  reviewer1,
		AssignedReviewerSecond: reviewer2,
		CreatedAt:              createdAt,
		MergedAt:               nil,
	}

	return pr, nil
}

func (p *postgresRepository) MergePullRequest(ctx context.Context, prID string) (*entity.PullRequest, error) {
	const query = `
UPDATE pull_requests 
SET 
status = 'MERGED', 
merged_at = CASE 
    WHEN status != 'MERGED' THEN NOW()
    ELSE merged_at
END
WHERE id = $1
RETURNING name, author_id, status, 
assigned_reviewer_first, assigned_reviewer_second,
created_at, merged_at
`
	pr := &entity.PullRequest{
		ID: prID,
	}

	err := p.db.QueryRow(ctx, query, prID).Scan(
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.AssignedReviewerFirst,
		&pr.AssignedReviewerSecond,
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &entity.PullRequest{}, entity.ErrPrNotFound
		}
		return &entity.PullRequest{}, err
	}

	return pr, nil
}

func (p *postgresRepository) ReassignReviewer(ctx context.Context, prID, userID string) (*entity.PullRequest, string, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return &entity.PullRequest{}, userID, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
		if err != nil {
			p.logger.Debug("Error: rollback")
		}
	}(tx, ctx)

	const query = `
SELECT name, author_id, status, 
assigned_reviewer_first, assigned_reviewer_second,
created_at, merged_at
FROM pull_requests 
WHERE id = $1
FOR UPDATE
`
	pr := &entity.PullRequest{
		ID: prID,
	}

	err = tx.QueryRow(ctx, query, prID).Scan(
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.AssignedReviewerFirst,
		&pr.AssignedReviewerSecond,
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &entity.PullRequest{}, userID, entity.ErrPrNotFound
		}
		return &entity.PullRequest{}, userID, err
	}
	if pr.Status == entity.MERGED {
		return &entity.PullRequest{}, userID, entity.ErrPrMerged
	}
	if pr.AssignedReviewerFirst != userID && pr.AssignedReviewerSecond != userID {
		return &entity.PullRequest{}, userID, entity.ErrUserNotAssigned
	}

	team, err := p.getTeamByUserID(ctx, tx, userID)
	if err != nil {
		return &entity.PullRequest{}, userID, err
	}

	newReviewer := findReplacementReviewer(team.Members, pr.AuthorID, userID, pr.AssignedReviewerFirst, pr.AssignedReviewerSecond)
	if newReviewer == "" {
		return &entity.PullRequest{}, userID, entity.ErrPrNoCandidate
	}

	const checkQuery = `
SELECT assigned_reviewer_first, assigned_reviewer_second 
FROM pull_requests 
WHERE id = $1
`
	var reviewer1, reviewer2 string
	err = tx.QueryRow(ctx, checkQuery, prID).Scan(&reviewer1, &reviewer2)
	if err != nil {
		return &entity.PullRequest{}, userID, err
	}

	var queryUpdate string
	if reviewer1 == userID {
		queryUpdate = `UPDATE pull_requests SET assigned_reviewer_first = $1 WHERE id = $2`
		pr.AssignedReviewerFirst = newReviewer
	} else if reviewer2 == userID {
		queryUpdate = `UPDATE pull_requests SET assigned_reviewer_second = $1 WHERE id = $2`
		pr.AssignedReviewerSecond = newReviewer
	}

	_, err = tx.Exec(ctx, queryUpdate, newReviewer, prID)
	if err != nil {
		return &entity.PullRequest{}, userID, err
	}

	if err = tx.Commit(ctx); err != nil {
		return &entity.PullRequest{}, userID, err
	}

	return pr, newReviewer, nil
}

func (p *postgresRepository) teamExists(ctx context.Context, tx pgx.Tx, teamName string) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)", teamName).Scan(&exists)
	return exists, err
}

func (p *postgresRepository) getTeamByUserID(ctx context.Context, tx pgx.Tx, userID string) (*entity.Team, error) {
	var teamName string

	const queryTeam = `
SELECT name FROM teams WHERE user_id = $1 LIMIT 1
`
	err := tx.QueryRow(ctx, queryTeam, userID).Scan(&teamName)

	if err != nil {
		return nil, entity.ErrUserNotFound
	}

	const queryUser = `
SELECT u.id, u.name, u.is_active 
FROM users u
JOIN teams t ON u.id = t.user_id
WHERE t.name = $1
`
	rows, err := tx.Query(ctx, queryUser, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	team := &entity.Team{
		Name: teamName,
	}

	for rows.Next() {
		var member entity.TeamMember
		err = rows.Scan(&member.ID, &member.Name, &member.IsActive)
		if err != nil {
			return nil, err
		}
		team.Members = append(team.Members, member)
	}

	return team, nil
}
