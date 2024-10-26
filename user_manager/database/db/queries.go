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

	updateBio = `
	UPDATE profiles SET bio = ? WHERE id = UNHEX(REPLACE(?, "-", ""));
	`

	createWorkspaceQuery = `
	INSERT INTO workspaces (id, name, description, user_id) VALUES (UNHEX(REPLACE(?, "-", "")), ?, ?, UNHEX(REPLACE(?, "-", "")));
	`

	createInvitationQuery = `
	INSERT INTO invitations (id, info, invitors_id, receiver_id, workspace_id) VALUES (UNHEX(REPLAC(?, "-", "")), ? , UNHEX(REPLACE(?, "-", "")), UNHEX(REPLACE(?, "-", "")), UNHEX(REPLACE(?, "-"; "")));
	`
	createRequestQuery = `
	INSERT INTO join_requests (id, info, reason, workspace_id, request_id) VALUES (UNHEX(REPLACE(?, "-", "")), ?, ?, UNHEX(REPLACE(?, "-", "")), UNHEX(REPLACE(?, "-", ""))),
	`

	updateName = `
	UPDATE contacts SET name = ? WHERE id = UNHEX(REPLACE(?, "-", ""))
	`

	joinRequestInfos = `
	SELECT 
		jri.id
		jri.info,
		jri.join_request_id
	FROM join_request_infos jri
	WHERE jri.user_id = UNHEX(REPLACE(?, "-", ""));
	`

	invitationInfos = `
	SELECT 
		i.id,
		i.info,
		i.invitation_id
	FROM invitation_infos i
	WHERE i.workspace_id = UNHEX(REPLACE(?, "-", ""));
	`

	completeWorkspace = `
	SELECT
		w.id,
		w.name,
		w.description,
		c.id,
		c.name,
		c.email,
		c.image_file_path,
		i.id,
		i.info,
		i.invitation_id
	FROM workspaces w
	LEFT JOIN workspaces_contacts wc ON w.id = wc.workspace_id
	LEFT JOIN contacts c ON wc.contact_id = c.id
	LEFT JOIN invitation_infos i ON w.id = i.workspace_id
	WHERE w.id = UNHEX(REPLACE(?, "-", ""));
	`
)
