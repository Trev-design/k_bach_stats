using RabbitMQ.Client;

namespace UserManagementSystem.Services.RabbitMQ;

public sealed class RabbitMQLoggingService(IMessageChannel channel) : RabbitMQBase<IMessageChannel>(channel), IHostedService, IAsyncDisposable
{
    public async ValueTask DisposeAsync()
    {
        await _channel.DisposeAsync();
        await _connection.DisposeAsync();
    }

    public async Task StartAsync(CancellationToken cancellationToken)
    {
        await StartBroker();
        await ComputeMessages();
    }

    public async Task StopAsync(CancellationToken cancellationToken)
    {
        await _channel.CloseAsync(cancellationToken);
        await _connection.CloseAsync(cancellationToken);
    }

    protected override Task ComputeMessages()
    {
        throw new NotImplementedException();
    }

    protected async override Task StartBroker()
    {
        var factory = new ConnectionFactory { Uri = new(URL) };
        _connection = await factory.CreateConnectionAsync();
        _channel = await _connection.CreateChannelAsync();

        await _channel.ExchangeDeclareAsync(Exchange, Kind, true, false, null, false);
        await _channel.QueueDeclareAsync(Queue, true, false, false, null, false);
        await _channel.QueueBindAsync(Queue, Exchange, RoutingKey, null, false);
    }
}