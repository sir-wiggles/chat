package main

// GetOrCreateUser is a query that will create a user if one does not exist or retrieve an
// existing one. There are three args in this query: name, email and avatar respectively.
const GetOrCreateUser = `
WITH row AS (
INSERT INTO
	users (name, email, avatar)
SELECT
	$1, $2, $3
WHERE
	NOT EXISTS (
		SELECT
			*
		FROM
			users
		WHERE
			name = $1
	) RETURNING *
)
SELECT
	id, uuid, name, email, avatar, FALSE as existing
FROM
	row
UNION
SELECT
	id, uuid, name, email, avatar, TRUE as existing
FROM
	users
WHERE
	name = $1;`
