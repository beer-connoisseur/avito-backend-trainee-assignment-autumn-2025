-- +goose Up
CREATE INDEX index_teams_user_id ON teams(user_id);

-- +goose Down
DROP INDEX index_teams_user_id;