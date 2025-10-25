using System.ComponentModel.DataAnnotations;

namespace UserManagementSystem.Models;

public class ImageUploadRequest
{
    [Required]
    public required string ContentType { get; set; }

    [Required]
    public required string FileName { get; set; }
} 