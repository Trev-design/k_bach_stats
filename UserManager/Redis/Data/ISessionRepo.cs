using UserManager.Redis.Models;

namespace UserManager.Redis.Data;

public interface ISessionRepo
{
    Task CreateSessionAsync(Session session);

    Task<Session?> GetSession(string sessionId);

    Task DeleteSessionAsync(string sessionId);
}