-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds (
    id serial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL,
    updated_at timestamp(0) with time zone NOT NULL,
    name text NOT NULL,
    url text UNIQUE NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feeds;
-- +goose StatementEnd
