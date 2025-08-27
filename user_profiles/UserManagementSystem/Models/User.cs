namespace UserManagementSystem.Models;

public class User
{
    public Guid Id { get; set; }
    public string Entity { get; set; } = null!;
    public Profile UserProfile { get; set; } = null!;
    public List<Contact> Contacts { get; set; } = [];
    public List<Workspace> Workspaces { get; set; } = [];
}

