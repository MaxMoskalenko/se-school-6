-- +goose Up
-- +goose StatementBegin
ALTER TABLE git_repositories ADD COLUMN last_checked_at TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE git_repositories DROP COLUMN last_checked_at;
-- +goose StatementEnd
