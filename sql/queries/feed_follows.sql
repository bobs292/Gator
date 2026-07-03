-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    values ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
JOIN users
    ON inserted_feed_follow.user_id = users.id
JOIN feeds
    ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM feed_follows
JOIN users
    ON feed_follows.user_id = users.id
JOIN feeds
    ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY feed_follows.created_at;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1
  AND feed_id = $2;