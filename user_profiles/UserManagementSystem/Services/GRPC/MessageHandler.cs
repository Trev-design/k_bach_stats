using System.Threading.Channels;
using Grpc.Core;
using UserManagementSystem.Grpc;
using UserManagementSystem.Services.Database;

namespace UserManagementSystem.Services.GRPC;

public class MessageHandler
{
    public RegistryRequest Request { init; private get; } = null!;
    public IServiceProvider ServiceProvider { init; private get; } = null!;
    public Channel<Response> MessagePipe { init; private get; } = null!;
    private readonly ulong index = 0;

    public async Task ComputeMessageAsync()
    {
        string message = "";
        try
        {
            using var scope = ServiceProvider.CreateScope();
            var dbContext = scope.ServiceProvider.GetRequiredService<AppDBContext>();
            await UserDBImpl.CreateUser(dbContext, Request.Name, Request.Email, Request.Entity);
            message = "ACCEPTED";
        }
        catch (Exception e)
        {
            Console.WriteLine(e.Message);
            message = "SOMETHING WENT WRONG";
        }
        finally
        {
            Response response = new()
            {
                Index = index,
                RegistryResponse = new()
                {
                    Status = message
                }
            };

            await MessagePipe.Writer.WriteAsync(response);
        }
    }
}