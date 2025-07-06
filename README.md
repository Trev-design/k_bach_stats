# A Deep Dive into Distributed Systems (In Progress)

## Dev Setup Almost Complete for Auth and Mailing

### Reminder
- Certificate files, key files, PEM files, and environment files are excluded from version control for security best practices.
- The docker-compose setup is not yet ready
- While the full development environment is still in progress, individual components are already covered by isolated test setups.

### Auth Responsibilities
The authentication service should:
1. Authorize users
2. Store sensitive user data securely
3. Store session data securely
4. Provide a fallback for password recovery
5. Implement role-based access control
6. Start user sessions
7. Refresh sessions
8. End sessions on user request
9. Handle session expiry and fallback mechanisms
10. Interact with the mailer service for verification processes

### Mailer Service Responsibilities
The mailer should:
1. Receive email payloads
2. Select the appropriate template based on the payload
3. Send emails securely

### Completed
- Component library setup in `auth`
- Component library tests in `auth`
- Component library setup in `mailer_server`

### In Progress / TODO
- Make TLS work
- Write better tests in auth service session
- Finalize component library in `mailer_server`
- Complete Docker Compose configuration
- Set up full application stack using Docker Compose
- Run initial data migration
- Create a full test environment for the distributed infrastructure