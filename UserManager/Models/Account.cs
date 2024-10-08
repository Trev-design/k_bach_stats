using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace UserManager.Models;

public class Account 
{
    [Key]
    [Column(TypeName = "binary(16)")]
    public Guid Id { get; set; } = Guid.NewGuid();

    [Required]
    public required string Entity { set; get; } 

    public User? AccountUser { get; set; }

    public ICollection<WorkSpace>? WorkSpaces { get; set; }

    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
}