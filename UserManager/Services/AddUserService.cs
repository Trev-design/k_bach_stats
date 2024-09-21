
using UserManager.Rabbit;

namespace UserManager.Services;

public class AddUserService(IServiceScopeFactory scopeFactory, RabbitConn conn) : IHostedService
{
    private readonly AddUserConsumer _consumer = new("account", "add_account", "add_account_request", conn, scopeFactory);

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
