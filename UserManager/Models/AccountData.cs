using System.Text.Json.Serialization;

namespace UserManager.Models;

public class AccountData
{
    [JsonPropertyName("entity")]
    public required string Entity { get; set; }

    [JsonPropertyName("username")]
    public required string Username { get; set;}

    [JsonPropertyName("email")]
    public required string Email { get; set;}
}