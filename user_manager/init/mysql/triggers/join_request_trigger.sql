CREATE TRIGGER increment_join_request_count
AFTER INSERT ON join_requests
FOR EACH ROW
BEGIN
    UPDATE workspaces
    SET join_request_count = join_request_count + 1
    WHERE id = NEW.workspace_id;
END;