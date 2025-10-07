using System.Threading.Channels;
using Grpc.Core;
using UserManagementSystem.Grpc;

namespace UserManagementSystem.Services.GRPC;

/// <summary>
/// handles all outcomming responses of the grpc servev
/// </summary>
public class ResponseHandler
{
    public IAsyncStreamWriter<RegistryResponse> Writer { init; private get; } = null!;
    public SemaphoreSlim Semaphore { init; private get; } = null!;
    public Channel<Response> MessagePipe { init; private get; } = null!;
    public CancellationTokenSource TokenSource { init; private get; } = null!;

    /// <summary>
    /// computes all messages and send it back in the right order
    /// </summary>
    /// <returns></returns>
    public async Task HandleResponseOutput()
    {
        Dictionary<ulong, RegistryResponse> responses = [];
        ulong next = 0;
        while (!TokenSource.Token.IsCancellationRequested)
        {
            try
            {
                Response response = await MessagePipe.Reader.ReadAsync();
                responses.Add(response.Index, response.RegistryResponse);

                while (true)
                {
                    bool ok = responses.TryGetValue(next, out RegistryResponse? currentResponse);
                    if (!ok) break;
                    if (currentResponse == null) throw new ArgumentNullException();
                    await Writer.WriteAsync(currentResponse);
                    responses.Remove(next);
                    Semaphore.Release();
                }
            }
            catch (NotSupportedException e)
            {
                Console.WriteLine(e.Message);
            }
            catch (ArgumentNullException e)
            {
                Console.WriteLine(e.Message);
            }
        }
    }
}