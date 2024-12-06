CREATE TABLE IF NOT EXISTS users(
    id BINARY(16) PRIMARY KEY,
    entity VARCHAR(255) NOT NULL UNIQUE,
    invitation_count INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS profiles(
    id BINARY(16) PRIMARY KEY,
    bio TEXT NOT NULL,
    user_id BINARY(16) NOT NULL UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS contacts(
    id BINARY(16) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    has_image VARCHAR(255) NOT NULL DEFAULT 'FALSE',
    profile_id BINARY(16) NOT NULL UNIQUE,
    user_id BINARY(16),
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS workspaces(
    id BINARY(16) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    join_request_count INT NOT NULL DEFAULT 0,
    user_id BINARY(16) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS workspaces_contacts(
    contact_id BINARY(16),
    workspace_id BINARY(16),
    PRIMARY KEY (contact_id, workspace_id),
    FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invitations(
    id BINARY(16) PRIMARY KEY,
    subject VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'UNREAD',
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    workspace_id BINARY(16) NOT NULL,
    user_id BINARY(16) NOT NULL,
    contact_id BINARY(16) NOT NULL,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS join_requests(
    id BINARY(16) PRIMARY KEY,
    subject VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'UNREAD',
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    profile_id BINARY(16) NOT NULL,
    workspace_id BINARY(16) NOT NULL,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
);