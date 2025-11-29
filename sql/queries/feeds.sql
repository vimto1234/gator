-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.id, feeds.created_at, feeds.updated_at, feeds.name, feeds.url, users.name AS username FROM feeds
JOIN users ON users.id = feeds.user_id;

-- name: FollowFeeds :one
WITH inserted_feed_follow AS (
    INSERT INTO user_feeds_link (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    ) 
    RETURNING *
)

SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users on users.id = inserted_feed_follow.user_id
INNER JOIN feeds on feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedbyURL :one
SELECT feeds.id, feeds.created_at, feeds.updated_at, feeds.name, feeds.url, users.name AS username FROM feeds
JOIN users ON users.id = feeds.user_id
WHERE url = $1;

-- name: GetAllFeedsUserFollows :many
SELECT users.name AS user_name, feeds.name AS feeds_name FROM user_feeds_link
INNER JOIN users on users.id = user_feeds_link.user_id
INNER JOIN feeds on feeds.id = user_feeds_link.feed_id
WHERE user_feeds_link.user_id = $1;

-- name: UnFollowByURLAndUsername :exec
DELETE FROM user_feeds_link 
WHERE user_feeds_link.user_id = $1 AND user_feeds_link.feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $1, last_fetched_at = $1
WHERE id = $2;

-- name: GetNextFeedToFetch :one
SELECT id, url FROM feeds
ORDER BY last_fetched_at NULLS FIRST;