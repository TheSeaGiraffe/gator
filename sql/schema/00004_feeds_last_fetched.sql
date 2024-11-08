-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds
ADD COLUMN last_fetched_at timestamp(0) with time zone;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds
DROP COLUMN last_fetched_at;
-- +goose StatementEnd
