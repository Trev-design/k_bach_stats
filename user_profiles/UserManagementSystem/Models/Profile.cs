using Microsoft.EntityFrameworkCore;

namespace UserManagementSystem.Models;

[Index(nameof(UserId), IsUnique = true)]
public class Profile
{
    public Guid Id { get; set; }
    public string? ImagePath { get; set; }
    public string? Description { get; set; }
    public Contact UserContact { get; set; } = null!;
    public Guid UserId { get; set; }
    public User ProfileUser { get; set; } = null!;
}


