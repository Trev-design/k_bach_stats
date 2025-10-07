using Grpc.Core;
using UserManagementSystem.Grpc;
using static UserManagementSystem.Grpc.UserRegistryService;

namespace UserManagementSystem.Services.GRPC;

/// <summary>
/// 
/// </summary>
/// <param name="provider"></param>
public class UserRegistryServiceImpl(IServiceProvider provider) : UserRegistryServiceBase
{
    private readonly IServiceProvider _serviceProvider = provider;

    public override async Task UserPrimaryStream(
        IAsyncStreamReader<RegistryRequest> requestStream,
        IServerStreamWriter<RegistryResponse> responseStream,
        ServerCallContext context)
    {
        await HandleStreamAsync(requestStream, responseStream, context);
    }

    public override async Task UserOverflowStream(
        IAsyncStreamReader<RegistryRequest> requestStream,
        IServerStreamWriter<RegistryResponse> responseStream,
        ServerCallContext context)
    {
        await HandleStreamAsync(requestStream, responseStream, context);
    }

    private async Task HandleStreamAsync(
        IAsyncStreamReader<RegistryRequest> requestStream,
        IServerStreamWriter<RegistryResponse> responseStream,
        ServerCallContext context)
    {
        StreamHandler handler = new()
        {
            RequestReader = requestStream,
            ResponseWriter = responseStream,
            Context = context,
            ServiceProvider = _serviceProvider
        };

        await handler.HandleMessageIncomeAsync();
    }
}