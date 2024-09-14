using System.ComponentModel.DataAnnotations;

namespace UserManager.Models;

public class Contact 
{
    [Key]
    public string Id { get; set; } = Guid.NewGuid().ToString();

    [Required]
    public required string Name { get; set; }

    [Required]
    public required string Email { get; set; }
}