using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace UserManager.Models;

public class WorkSpace
{
    [Key]
    [Column(TypeName = "binary(16)")]
    public Guid Id { get; set;} = Guid.NewGuid();

    [Required]
    public required string Name { get; set;}

    [Required]
    public required Guid AccountId { get; set;}

    public Account? Account { get; set;}

    public ICollection<Contact>? Contacts { get; set;}

    public DateTime CreatedAt { get; set;} = DateTime.UtcNow;
    public DateTime UpdatedAt { get; set;} = DateTime.UtcNow;
}