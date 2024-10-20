-- name: CreateFeedFollow :one
WITH new_record AS (
    INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
    VALUES($1, $2, $3, $4)
    RETURNING *
)
SELECT
    new_record.id,
    new_record.created_at,
    new_record.updated_at,
    users.name AS user_name,
    feeds.name AS feed_name
FROM new_record
INNER JOIN users ON new_record.user_id = users.id
INNER JOIN feeds ON new_record.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT
    users.name AS user_name,
    feeds.name AS feed_name
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;
