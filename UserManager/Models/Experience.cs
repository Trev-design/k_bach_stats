using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace UserManager.Models;

public class Experience
{
    [Key]
    [Column(TypeName = "binary(16)")]
    public Guid Id { get; set; }

    [Required]
    public required string Tag { get; set; }

    public ICollection<Contact>? Contacts { get; set; }

    public ICollection<SelfAssessment>? SelfAssessments { get; set; }

    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
}