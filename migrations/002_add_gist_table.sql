-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gists (
    file_id text primary key,
    user_id text not null references users(id),
    gist_title text not null,
    forked_from text references users(id),
    short_url text not null,
    view_count numeric not null default 0,
    is_public boolean not null default true,
    is_deleted boolean not null default false,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gists;
-- +goose StatementEnd
