CREATE TRIGGER increment_invitation_count
AFTER INSERT ON invitations
FOR EACH ROW
BEGIN
    UPDATE users
    SET invitation_count = invitation_count + 1
    WHERE id = NEW.user_id;
END;