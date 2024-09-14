using System.ComponentModel.DataAnnotations;

namespace UserManager.Models;

public class User
{
    [Key]
    public string Id { get; set; } = new Guid().ToString();

    [Required]
    public required string Name { get; set; }

    public List<User>? Friends { get; set; }
}