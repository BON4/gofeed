-- name: CreatePost :one
INSERT INTO Posts (
       content,
       posted_on,
       posted_by,
       score
) VALUES (
  $1, $2, $3, $4
) RETURNING post_id;

-- name: DeletePost :exec
DELETE FROM Posts WHERE post_id = $1;

-- name: ListPosts :many
SELECT * FROM Posts offset $2 limit $1;
