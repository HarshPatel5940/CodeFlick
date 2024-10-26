-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS replies (
    id text primary key,
    user_id text not null references users(id),
    gist_id text not null references gists(file_id),
    message text not null,
    is_deleted boolean not null default false,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE replies;
-- +goose StatementEnd
