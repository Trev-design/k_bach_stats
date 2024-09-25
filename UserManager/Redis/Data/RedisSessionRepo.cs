using System.Text.Json;
using StackExchange.Redis;
using UserManager.Redis.Models;

namespace UserManager.Redis.Data;

public class RedisSessionRepo : ISessionRepo
{
    private readonly IConnectionMultiplexer _redis;

    public RedisSessionRepo(IConnectionMultiplexer redis)
    {
        Console.WriteLine("starting redis session repository");
        _redis = redis;
    }

    public async Task CreateSessionAsync(Session session)
    {
        ArgumentNullException.ThrowIfNull(session);

        var db = _redis.GetDatabase(2);
        var serialSession = JsonSerializer.Serialize(session);
        var expiryTime = DateTimeOffset.Now.AddSeconds(60 * 60 *24);
        var expiry = expiryTime.DateTime.Subtract(DateTime.Now);
        await db.StringSetAsync(session.Id, serialSession, expiry);
    }

    public async Task DeleteSessionAsync(string sessionId)
    {
        if (sessionId == null)
        {
            throw new ArgumentNullException(sessionId);
        }

        var db = _redis.GetDatabase(2);
        await db.KeyDeleteAsync(sessionId);
    }

    public async Task<Session?> GetSession(string sessionId)
    {
        if (sessionId == null)
        {
            throw new ArgumentNullException(sessionId);    
        }

        var db = _redis.GetDatabase(2);
        string? sessionString = await db.StringGetAsync(sessionId);

        if (sessionString == null)
        {
            return null;
        }

        return JsonSerializer.Deserialize<Session>(sessionString);
    }
}

