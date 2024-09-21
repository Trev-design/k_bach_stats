
using UserManager.Rabbit;

namespace UserManager.Services;

public class StopSessionService(IServiceScopeFactory scopeFactory, RabbitConn conn) : IHostedService
{
    private readonly StopSessionConsumer _consumer = new("session", "stop_user_session", "remove_user_session", conn, scopeFactory);
    public Task StartAsync(CancellationToken cancellationToken)
    {
        _consumer.Consume(cancellationToken);
        return Task.CompletedTask;
    }

    public Task StopAsync(CancellationToken cancellationToken)
    {
        _consumer.StopConsumer();
        return Task.CompletedTask;
    }
}
