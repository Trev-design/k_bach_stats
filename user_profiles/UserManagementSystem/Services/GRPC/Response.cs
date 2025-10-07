using UserManagementSystem.Grpc;

namespace UserManagementSystem.Services.GRPC;

/// <summary>
/// a data class to compute responses with ease
/// </summary>
public class Response
{
    public ulong Index { get; set; }
    public RegistryResponse RegistryResponse { get; set; } = null!;
}