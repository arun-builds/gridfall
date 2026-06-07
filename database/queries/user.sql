-- name: GetUser :one
select * from users
where id = $1;

-- name: ListUsers :many
select * from users
order by name;

-- name: CreateAuthor :one
insert into users(
name, type
) values (
$1, $2
)
returning *;

-- name: UpdateUserName :exec
UPDATE users
  set name = $2
WHERE id = $1;

-- name: UpdateUserType :exec
UPDATE users
  set type = $2
WHERE id = $1;

-- name: DeleteAuthor :exec
DELETE FROM users
WHERE id = $1;
