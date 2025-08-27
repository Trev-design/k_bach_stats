using Microsoft.EntityFrameworkCore;

namespace UserManagementSystem.Models;

[Index(nameof(ProfileId), IsUnique = true)]
[Index(nameof(Email), IsUnique = true)]
public class Contact
{
    public Guid Id { get; set; }
    public string Name { get; set; } = null!;
    public string Email { get; set; } = null!;
    public Guid ProfileId { get; set; }
    public Profile UserProfile { get; set; } = null!;
    public List<User> Users { get; set; } = [];
    public List<Workspace> Workspaces { get; set; } = [];
}

