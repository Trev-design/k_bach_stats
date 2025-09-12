using System.Threading.Channels;
using Grpc.Core;
using UserManagementSystem.Grpc;

namespace UserManagementSystem.Services.GRPC;

public class StreamHandler
{
    public IAsyncStreamReader<RegistryRequest> RequestReader { init; private get; } = null!;
    public IAsyncStreamWriter<RegistryResponse> ResponseWriter { init; private get; } = null!;
    public ServerCallContext Context { init; private get; } = null!;
    public IServiceProvider ServiceProvider { init; private get; } = null!;
    private readonly Channel<Response> MessageStream = Channel.CreateUnbounded<Response>();
    private readonly SemaphoreSlim Semaphore = new(10);

    public async Task HandleMessageIncomeAsync()
    {
        using var cancelationTokenSource = CancellationTokenSource.CreateLinkedTokenSource(Context.CancellationToken);

        var responseTask = Task.Run(async () =>
        {
            ResponseHandler handler = new()
            {
                Writer = ResponseWriter,
                Semaphore = Semaphore,
                MessagePipe = MessageStream,
                TokenSource = cancelationTokenSource
            };

            await handler.HandleResponseOutput();
        });

        try
        {
            await foreach (var request in RequestReader.ReadAllAsync(cancelationTokenSource.Token))
            {
                await Semaphore.WaitAsync();
                
                _ = Task.Run(async () =>
                {
                    MessageHandler handler = new()
                    {
                        Request = request,
                        ServiceProvider = ServiceProvider,
                        MessagePipe = MessageStream
                    };

                    await handler.ComputeMessageAsync();
                });
            }
        }
        finally
        {
            await responseTask;
        }

    } 
}