namespace UserManagementSystem.Models;

public class Workspace
{
    public Guid Id { get; set; }
    public string Name { get; set; } = null!;
    public Guid UserId { get; set; }
    public User User { get; set; } = null!;
    public List<Contact> Contacts { get; set; } = [];
    public List<ChatRoom> ChatRooms { get; set; } = [];
}

