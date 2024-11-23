package sqldb

const (
	userIDQuery = `
	SELECT u.id FROM users u WHERE u.entity = ?;
	`

	userByIDQuery = `
	SELECT 
		p.id,
		p.bio,
		c.id,
		c.name,
		c.email,
		c.image_file_path
	FROM users u
	JOIN profiles p ON u.id = p.user_id 
	JOIN contacts c ON p.id = c.profile_id 
	WHERE u.id = UNHEX(REPLACE(?, "-", ""))
	`
)
