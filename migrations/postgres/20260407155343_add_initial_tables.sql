-- +goose Up
-- +goose StatementBegin
-- users table to store user information, such as email
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- repositories table to store information about git repositories, name and owner
-- are separated to allow easier querying, last_seen_tag indicates the latest tag seen in the repository, 
-- which is used to determine if there are new releases
CREATE TABLE IF NOT EXISTS git_repositories (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner VARCHAR(255) NOT NULL,
    last_seen_tag VARCHAR(255),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- unique case insensative index on owner and name to prevent duplicate repositories with different cases
CREATE UNIQUE INDEX idx_git_repositories_owner_name_lower ON git_repositories (LOWER(owner), LOWER(name));

-- repository_subscriptions table to store user subscriptions to repositories, 
-- confirmed_at and unsubscribe are used to track the subscription status, 
-- allowing for soft deletes and re-subscriptions without losing history
CREATE TABLE IF NOT EXISTS repository_subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    repository_id UUID NOT NULL,
    confirmed_at TIMESTAMPTZ,
    unsubscribed_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (repository_id) REFERENCES git_repositories (id),
    
    UNIQUE (user_id, repository_id)
);

-- doi_subscription_tokens table to store tokens for confirming subscribing and unsubscribing
-- action field is "not-strict" enum, allowing to support multiple actions without adding new migration files
CREATE TABLE IF NOT EXISTS doi_subscription_tokens (
    id UUID PRIMARY KEY, -- uuid is confirmation/unsubscribe token
    subscription_id UUID NOT NULL,
    action SMALLINT NOT NULL, -- 0 for subscribe, 1 for unsubscribe

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (subscription_id) REFERENCES repository_subscriptions (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS doi_subscription_tokens;
DROP TABLE IF EXISTS repository_subscriptions;
DROP INDEX IF EXISTS idx_git_repositories_owner_name_lower;
DROP TABLE IF EXISTS git_repositories;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd