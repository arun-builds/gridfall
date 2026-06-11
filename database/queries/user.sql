-- name: GetUser :one
select * from users
where id = $1;

-- name: GetUserByEmail :one
select * from users
where email = $1;

-- name: ListUsers :many
select * from users
order by name;

-- name: CreateRegisteredUser :one
insert into users(
  name, email, password_hash, type
) values (
  $1, $2, $3, 'registered'
)
returning *;

-- name: CreateGuestUser :one
insert into users(
  name, type
) values (
  $1, 'guest'
)
returning *;

-- name: UpgradeGuestToRegistered :one
UPDATE users
  SET email = $2,
      password_hash = $3,
      type = 'registered'
WHERE id = $1 AND type = 'guest'
RETURNING *;

-- name: UpdateUserName :exec
UPDATE users
  set name = $2
WHERE id = $1;

-- name: UpdateUserType :exec
UPDATE users
  set type = $2
WHERE id = $1;

-- name: UpdateUserRole :exec
UPDATE users
  set role = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
