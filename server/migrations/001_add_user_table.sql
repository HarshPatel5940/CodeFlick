-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id text primary key,
    name text not null,
    email text not null UNIQUE,
    auth_provider text not null,
    is_admin boolean not null default false,
    is_premium boolean not null default false,
    is_deleted boolean not null default false,
    created_at timestamp  not null default now(),
    updated_at timestamp  not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
