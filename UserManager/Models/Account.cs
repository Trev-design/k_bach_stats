using System.ComponentModel.DataAnnotations;

namespace UserManager.Models;
public class Account
{
    [Key]
    public string Id { get; set; } = new Guid().ToString();

    [Required]
    public required string Entity { get; set; }

    [Required]
    public required Contact Contact { get; set; }

    [Required]
    public required User User { get; set; }

    public List<Workspace>? Workspaces { get; set; }
}