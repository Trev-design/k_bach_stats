using UserManagementSystem.Grpc;

namespace UserManagementSystem.Services.GRPC;

public class Response
{
    public ulong Index { get; set; }
    public RegistryResponse RegistryResponse { get; set; } = null!;
}