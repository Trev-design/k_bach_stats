
using UserManager.Rabbit;

namespace UserManager.Services;

public class DeleteUserService(IServiceScopeFactory scopeFactory, RabbitConn conn) : IHostedService
{
    private readonly DeleteUserConsumer _consumer = new("account", "delete_account", "delete_account_request", conn, scopeFactory);

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
