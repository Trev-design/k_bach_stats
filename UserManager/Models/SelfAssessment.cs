using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace UserManager.Models;

public class SelfAssessment
{
    [Key]
    [Column(TypeName = "binary(16)")]
    public Guid Id { get; set; } = Guid.NewGuid();

    [Required]
    public required int Evaluation { get; set; }

    public ICollection<Experience>? Experiences { get; set; }
}