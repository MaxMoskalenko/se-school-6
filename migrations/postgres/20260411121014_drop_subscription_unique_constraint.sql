-- +goose Up
-- +goose StatementBegin
-- removing unique constraint on subscriptions, but keeping index for search
ALTER TABLE repository_subscriptions DROP CONSTRAINT repository_subscriptions_user_id_repository_id_key;
CREATE INDEX idx_repository_subscriptions_user_repo ON repository_subscriptions (user_id, repository_id);

-- replace expression-based unique index with a plain unique constraint
-- owner/name are now normalized to lowercase at the application level
DROP INDEX idx_git_repositories_owner_name_lower;
ALTER TABLE git_repositories ADD CONSTRAINT uq_git_repositories_owner_name UNIQUE (owner, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE git_repositories DROP CONSTRAINT uq_git_repositories_owner_name;
CREATE UNIQUE INDEX idx_git_repositories_owner_name_lower ON git_repositories (LOWER(owner), LOWER(name));

DROP INDEX idx_repository_subscriptions_user_repo;
ALTER TABLE repository_subscriptions ADD CONSTRAINT repository_subscriptions_user_id_repository_id_key UNIQUE (user_id, repository_id);
-- +goose StatementEnd
