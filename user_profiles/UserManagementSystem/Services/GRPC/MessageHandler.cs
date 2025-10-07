using System.Threading.Channels;
using Grpc.Core;
using Microsoft.EntityFrameworkCore;
using MySqlConnector;
using UserManagementSystem.Grpc;
using UserManagementSystem.Services.Database;

namespace UserManagementSystem.Services.GRPC;

/// <summary>
/// class to compute incomming requests.
///  
/// spawns on every incomming request
/// </summary>
public class MessageHandler
{
    public RegistryRequest Request { init; private get; } = null!;
    public IServiceProvider ServiceProvider { init; private get; } = null!;
    public Channel<Response> MessagePipe { init; private get; } = null!;
    private ulong index = 0;

    /// <summary>
    /// computes incomming requests
    /// </summary>
    /// <returns></returns>
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
        catch (DbUpdateException ex) when (IsUniqueConstraintViolation(ex))
        {
            message = "DUPLICATE";
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
            index++;
        }
    }

    private static bool IsUniqueConstraintViolation(DbUpdateException ex)
    {
        if (ex.InnerException is MySqlException mysqlEx)
        {
            // MySQL error code 1062 => Duplicate entry for key (unique constraint violation)
            return mysqlEx.Number == 1062;
        }

        return false;
    }
}