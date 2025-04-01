-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gists (
    file_id text primary key,
    user_id text not null references users(id),
    gist_title text not null,
    forked_from text references users(id),
    short_url text not null UNIQUE,
    file_name text not null,
    view_count numeric not null default 0,
    is_public boolean not null default true,
    is_deleted boolean not null default false,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

-- Create the trigger function
CREATE OR REPLACE FUNCTION set_forked_from_default()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if forked_from is NULL and set it to user_id
    IF NEW.forked_from IS NULL THEN
        NEW.forked_from := NEW.user_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create the trigger
CREATE TRIGGER set_forked_from_before_insert
BEFORE INSERT ON gists
FOR EACH ROW
EXECUTE FUNCTION set_forked_from_default();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gists;
-- +goose StatementEnd
