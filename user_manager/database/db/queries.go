package db

const (
	insertNewUser = `INSERT INTO users (id, entity) VALUES (UNHEX(REPLACE(?, "-", "")), ?);`

	insertNewProfile = `INSERT INTO profiles (id, bio, user_id) VALUES (UNHEX(REPLACE(?, "-", "")), ?, UNHEX(REPLACE(?, "-", "")));`

	insertNewContact = `INSERT INTO contacts (id, name, email, image_file_path, profile_id) VALUES (UNHEX(REPLACE(?, "-", "")), ?, ?, ?, UNHEX(REPLACE(?, "-", "")));`

	createDatatabas = `CREATE DATABASE IF NOT EXISTS user_database;`

	insertRating = `INSERT INTO self_assessments (id, rating) VALUES (UNHEX(REPLACE(?, "-", "")), ?);`

	dropDatabase = `DROP DATABASE user_database;`

	removeUser = `DELETE FROM users WHERE entity = ?`

	selectUser = `
	SELECT
		u.id,
		u.entity,
		u.requests,
		p.id, 
		p.bio,
		c.id,
		c.name,
		c.email,
		c.image_file_path,
		w.id,
		w.name
	FROM users u
	JOIN profiles p ON u.id = p.user_id
	JOIN contacts c ON p.id = c.profile_id
	LEFT JOIN workspaces w ON u.id = w.user_id 
	WHERE u.id = UNHEX(REPLACE(?, "-", ""));
	`

	userCredentials = `
	SELECT u.id FROM users u WHERE u.entity = ?;
	`
)
