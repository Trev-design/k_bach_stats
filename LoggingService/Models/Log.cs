using System.ComponentModel.DataAnnotations;

namespace LoggingService.Models;

public class Log 
{
    [Key]
    public Guid Id {get; set;}
    public string? DescriptionHeader {get; set;}
    public string? LogMessage {get; set;}
    public string? Error {get; set;}
}