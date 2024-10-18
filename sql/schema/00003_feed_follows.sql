-- +goose Up
-- +goose StatementBegin
CREATE TABLE feed_follows (
    id serial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL,
    updated_at timestamp(0) with time zone NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id int NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feed_follows;
-- +goose StatementEnd
