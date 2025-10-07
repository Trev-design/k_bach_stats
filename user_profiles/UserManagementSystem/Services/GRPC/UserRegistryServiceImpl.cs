using Grpc.Core;
using UserManagementSystem.Grpc;
using static UserManagementSystem.Grpc.UserRegistryService;

namespace UserManagementSystem.Services.GRPC;

/// <summary>
/// starting the streams of the grpc server
/// </summary>
/// <param name="provider"></param>
public class UserRegistryServiceImpl(IServiceProvider provider) : UserRegistryServiceBase
{
    private readonly IServiceProvider _serviceProvider = provider;

    /// <summary>
    /// starting the primary stream
    /// </summary>
    /// <param name="requestStream"></param>
    /// <param name="responseStream"></param>
    /// <param name="context"></param>
    /// <returns></returns>
    public override async Task UserPrimaryStream(
        IAsyncStreamReader<RegistryRequest> requestStream,
        IServerStreamWriter<RegistryResponse> responseStream,
        ServerCallContext context)
    {
        await HandleStreamAsync(requestStream, responseStream, context);
    }

    /// <summary>
    /// starting the secondary stream
    /// </summary>
    /// <param name="requestStream"></param>
    /// <param name="responseStream"></param>
    /// <param name="context"></param>
    /// <returns></returns>
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