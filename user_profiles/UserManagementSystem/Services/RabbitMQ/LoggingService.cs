using RabbitMQ.Client;

namespace UserManagementSystem.Services.RabbitMQ;

/// <summary>
/// 
/// </summary>
/// <param name="channel"></param>
public sealed class RabbitMQLoggingService(IMessageChannel channel) : RabbitMQBase<IMessageChannel>(channel), IHostedService, IAsyncDisposable
{
    private Task _messageTask = null!;
    public async ValueTask DisposeAsync()
    {
        await _channel.DisposeAsync();
        await _connection.DisposeAsync();
    }

    public async Task StartAsync(CancellationToken cancellationToken)
    {
        await StartBroker();
        _messageTask = Task.Run(async () =>
        {
            await ComputeMessages();
        }, cancellationToken);
    }

    public async Task StopAsync(CancellationToken cancellationToken)
    {
        _messageChannel.Complete();
        await _messageTask;
        await _channel.CloseAsync(cancellationToken);
        await _connection.CloseAsync(cancellationToken);
    }

    protected async override Task ComputeMessages()
    {
        await foreach (var message in _messageChannel.GetMessagePipe())
        {
            await _channel.BasicPublishAsync(Exchange, RoutingKey, message);
        }
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