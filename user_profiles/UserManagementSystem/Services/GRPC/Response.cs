using UserManagementSystem.Grpc;

namespace UserManagementSystem.Services.GRPC;

/// <summary>
/// 
/// </summary>
public class Response
{
    public ulong Index { get; set; }
    public RegistryResponse RegistryResponse { get; set; } = null!;
}