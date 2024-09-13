using System.ComponentModel.DataAnnotations;

namespace UserManager.Redis.Models;

public class Session 
{
    [Required]
    public required string ID { get; set; }

    [Required]
    public required string Name { get; set; }

    [Required]
    public required string Account { get; set; }

}