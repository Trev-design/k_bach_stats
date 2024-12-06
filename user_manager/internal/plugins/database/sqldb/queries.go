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
		c.has_image
	FROM users u
	JOIN profiles p ON u.id = p.user_id 
	JOIN contacts c ON p.id = c.profile_id 
	WHERE u.id = UNHEX(REPLACE(?, "-", ""))
	`

	insertInvitationQuery = `INSERT INTO invitations (id, subject, message, workspace_id, user_id, contact_id) VALUES (UNHEX(REPLACE(?, "-", "")), ?, ?, UNHEX(REPLACE(?, "-", "")), UNHEX(REPLACE(?, "-", "")), UNHEX(REPLACE(?, "-", "")));`

	insertJoinRequestQuery = `INSERT INTO join_requests (id, subject, message, profile_id, workspace_id) VALUES (UNHEX(REPLACE(?, "-", "")), ?, ?, UNHEX(REPLACE(?, "-", "")), UNHEX(REPLACE(?, "-", "")));`

	selectInvitationsQuery = `SELECT i.id, i.subject, i.message, i.workspace_id, i.user_id, i.contact_id FROM invitations i WHERE i.sent_at <= CURDATE() - INTERVAL 30 DAY;`

	selectJoinRequestsQuery = `SELECT j.id, j.subject, j.message, j.profile_id, j.workspace_id FROM join_requests j WHERE j.sent_at <= CURDATE() - INTERVAL 30 DAY;`

	insertWorkspaceQuery = "INSERT INTO workspaces (id, name, description, user_id) VALUES (UNHEX(REPLACE(?, '-', '')), ?, ?, UNHEX(REPLACE(?, '-', '')));"
)
