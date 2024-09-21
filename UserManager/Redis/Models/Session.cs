using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace UserManager.Redis.Models;

public class Session 
{
    [Required]
    [JsonPropertyName("id")]
    public required string Id { get; set; }

    [Required]
    [JsonPropertyName("name")]
    public required string Name { get; set; }

    [Required]
    [JsonPropertyName("account")]
    public required string Account { get; set; }

}