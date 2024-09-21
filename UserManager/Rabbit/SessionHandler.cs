using System.Text;
using System.Text.Json;
using UserManager.Redis.Data;
using UserManager.Redis.Models;

namespace UserManager.Rabbit;

public class SessionHandler(IServiceScopeFactory scopeFactory)
{
    private readonly IServiceScopeFactory _scopeFactory = scopeFactory;

    public async Task StartSession(ReadOnlyMemory<byte> body)
    {
        Console.WriteLine("IAmReadyToStart");

        using var scope = _scopeFactory.CreateAsyncScope();
        var context = scope.ServiceProvider.GetRequiredService<ISessionRepo>();

        Console.WriteLine("IAmReadyToStart");
        var payload = body.ToArray();
        Console.WriteLine(payload.Length);
        string jsonData = Encoding.UTF8.GetString(payload);
        Console.WriteLine(jsonData);
        Session session = JsonSerializer.Deserialize<Session>(jsonData) ?? throw new ArgumentException("invalid credentials");

        Console.WriteLine("start session");

        await context.CreateSessionAsync(session);
    }

    public async Task StopSession(ReadOnlyMemory<byte> body)
    {
        using var scope = _scopeFactory.CreateAsyncScope();
        var context = scope.ServiceProvider.GetRequiredService<ISessionRepo>();

        string jsonData = Encoding.UTF8.GetString(body.ToArray());
        Session session = JsonSerializer.Deserialize<Session>(jsonData) ?? throw new ArgumentException("invalid credentials");

        Session sessionToDelete = await context.GetSession(session.Id) ?? throw new ArgumentException("invalid session");

        await context.DeleteSessionAsync(sessionToDelete.Id.ToString());
    }
}


