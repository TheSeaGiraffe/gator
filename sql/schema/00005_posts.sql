-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
    id serial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL,
    updated_at timestamp(0) with time zone NOT NULL,
    title text NOT NULL,
    url text UNIQUE NOT NULL,
    description text,
    published_at timestamp(0) with time zone NOT NULL,
    feed_id int NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd
