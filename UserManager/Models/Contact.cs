using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace UserManager.Models;

public class Contact
{
    [Key]
    [Column(TypeName = "binary(16)")]
    public Guid Id { get; set; } = Guid.NewGuid();

    [Required]
    public required string Name { get; set; }

    [Required]
    public required string Email { get; set; }

    public string ProfileImagepath { get; set; } = string.Empty;

    [Required]
    [Column(TypeName = "binary(16)")]
    public required Guid ProfileId { get; set; }

    public Profile? Profile { get; set; }

    public ICollection<User>? Users { get; set; }
    public ICollection<WorkSpace>? WorkSpaces { get; set;}
    public ICollection<Experience>? Experiences { get; set; }

    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
    public DateTime UpdatedAt { get; set; } = DateTime.UtcNow;
}