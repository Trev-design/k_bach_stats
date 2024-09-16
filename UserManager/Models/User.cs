using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace UserManager.Models;

public class User 
{
    [Key]
    [Column(TypeName = "binary(16)")]
    public Guid Id { get; set; } = Guid.NewGuid();

    [Required]
    [Column(TypeName = "binary(16)")]
    public required Guid AccountId { get; set; }

    public Account? UserAccount { get; set; }

    public Profile? Profile { get; set; }

    public ICollection<Contact>? Contacts { get; set; }

    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;
}