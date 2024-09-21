
using UserManager.Rabbit;

namespace UserManager.Services;

public class StartSessionService(IServiceScopeFactory scopeFactory, RabbitConn conn) : IHostedService
{
    private readonly StartSessionConsumer _consumer = new("session", "start_user_session", "send_session_credentials", conn, scopeFactory);

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


