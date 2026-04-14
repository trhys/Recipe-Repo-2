-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (id, created_at, expires_at, user_id)
VALUES (
	$1,
	NOW(),
	NOW() + interval '30 days',
	$2
) RETURNING *;

-- name: GetRefreshToken :one
SELECT user_id FROM refresh_tokens
WHERE id = $1 AND expires_at > NOW() AND revoked_at IS NULL;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE id = $1;
