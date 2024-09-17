using UserManager.Data;
using UserManager.Rabbit;
using UserManager.Redis.Data;

namespace UserManager.Serices;

public class RabbitConsumerService(IServiceScopeFactory scopeFactory) : IHostedService
{
    private readonly Consumer _consumer = new(scopeFactory);

    public Task StartAsync(CancellationToken cancellationToken)
    {
        _ = _consumer.ConsumeSessionRequests(cancellationToken);
        return Task.CompletedTask;
    }

    public Task StopAsync(CancellationToken cancellationToken)
    {
        _consumer.StopConsumer();
        return Task.CompletedTask;
    }
}


