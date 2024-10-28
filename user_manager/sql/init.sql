CREATE TABLE IF NOT EXISTS users(
    id BINARY(16) PRIMARY KEY,
    entity VARCHAR(255) NOT NULL UNIQUE,
    requests BIGINT UNSIGNED NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS profiles(
    id BINARY(16) PRIMARY KEY,
    bio TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    user_id BINARY(16) NOT NULL UNIQUE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS contacts(
    id BINARY(16) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    image_file_path VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    profile_id BINARY(16) NOT NULL UNIQUE,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users_contacts(
    user_id BINARY(16),
    contact_id BINARY(16),
    PRIMARY KEY (user_id, contact_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS workspaces(
    id BINARY(16) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    requests BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    user_id BINARY(16) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS workspaces_contacts(
    workspace_id BINARY(16),
    contact_id BINARY(16),
    PRIMARY KEY (workspace_id, contact_id),
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
    FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS experiences(
    id BINARY(16) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    user_id BINARY(16),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS ratings(
    id BINARY(16) PRIMARY KEY,
    rating INT NOT NULL,
    experience_id BINARY(16),
    user_id BINARY(16),
    FOREIGN KEY (experience_id) REFERENCES experiences(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS experiences_profiles(
    experience_id BINARY(16),
    profile_id BINARY(16),
    PRIMARY KEY (experience_id, profile_id),
    FOREIGN KEY (experience_id) REFERENCES experiences(id) ON DELETE CASCADE,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invitations(
    id BINARY(16) PRIMARY KEY,
    info TEXT NOT NULL,
    ignored TINYINT(1) NOT NULL DEFAULT 0,
    invitors_id BINARY(16) NOT NULL,
    receiver_id BINARY(16) NOT NULL,
    workspace_id BINARY(16) NOT NULL,
    FOREIGN KEY (invitors_id) REFERENCES contacts(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS join_requests(
    id BINARY(16) PRIMARY KEY,
    info TEXT NOT NULL,
    reason TEXT NOT NULL,
    ignored TINYINT(1) NOT NULL DEFAULT 0,
    workspace_id BINARY(16) NOT NULL,
    request_id BINARY(16) NOT NULL,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
    FOREIGN KEY (request_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS invitation_infos(
    id BINARY(16) PRIMARY KEY,
    info TEXT NOT NULL,
    is_read TINYINT(1) NOT NULL DEFAULT 0,
    ignored TINYINT(1) NOT NULL DEFAULT 0,
    workspace_id BINARY(16) NOT NULL,
    invitations_id BINARY(16) NOT NULL,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
    FOREIGN KEY (invitations_id) REFERENCES invitations(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS join_request_infos(
    id BINARY(16) PRIMARY KEY,
    info TEXT NOT NULL,
    is_read TINYINT(1) NOT NULL DEFAULT 0,
    ignored TINYINT(1) NOT NULL DEFAULT 0,
    user_id BINARY(16) NOT NULL,
    join_request_id BINARY(16) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (join_request_id) REFERENCES join_requests(id) ON DELETE CASCADE
);
