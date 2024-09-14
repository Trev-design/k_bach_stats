using System.ComponentModel.DataAnnotations;

namespace UserManager.Models;

public class Workspace
{
    [Key]
    public string Id { get; set; } = new Guid().ToString();

    [Required]
    public required Contact Admin { get; set; }
    
    public List<Contact>? Collaborators { get; set; }
}