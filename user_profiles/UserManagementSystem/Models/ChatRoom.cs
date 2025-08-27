namespace UserManagementSystem.Models;

public class ChatRoom
{
    public Guid Id { get; set; }
    public string Topic { get; set; } = null!;
    public string Reference { get; set; } = null!;
    public Guid WorkspaceId { get; set; }
    public Workspace Workspace { get; set; } = null!;
    public List<Contact> Contacts { get; set; } = [];
}

