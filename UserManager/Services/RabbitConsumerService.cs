using UserManager.Rabbit;
using UserManager.Redis.Data;

namespace UserManager.Serices;

public class RabbitConsumerService(ISessionRepo repo) : IHostedService
{
    private readonly Consumer _consumer = new(repo);

    public async Task StartAsync(CancellationToken cancellationToken)
    {
        await _consumer.ConsumeSessionRequests(cancellationToken);
    }

    public Task StopAsync(CancellationToken cancellationToken)
    {
        _consumer.StopConsumer();
        return Task.CompletedTask;
    }
}


