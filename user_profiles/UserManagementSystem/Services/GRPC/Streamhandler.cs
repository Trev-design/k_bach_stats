using Grpc.Core;
using UserManagementSystem.Grpc;

namespace UserManagementSystem.Services.GRPC;

public class StreamHandler
{
    public IAsyncStreamReader<RegistryRequest> RequestReader { init; private get; } = null!;
    public IAsyncStreamWriter<RegistryResponse> ResponseWriter { init; private get; } = null!;
    public ServerCallContext Context { init; private get; } = null!;
    public IServiceProvider ServiceProvider { init; private get; } = null!;

    public async Task HandleStreamAsync()
    {
        RequestHandler handler = new()
        {
            RequestReader = RequestReader,
            ResponseWriter = ResponseWriter,
            Context = Context,
            ServiceProvider = ServiceProvider
        };

        await handler.HandleMessageIncomeAsync();
    }
}