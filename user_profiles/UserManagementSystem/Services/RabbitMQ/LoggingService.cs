using RabbitMQ.Client;

namespace UserManagementSystem.Services.RabbitMQ;

/// <summary>
/// the rabbitmq implementation
/// </summary>
/// <param name="channel"></param>
public sealed class RabbitMQLoggingService(IMessageChannel channel) : RabbitMQBase<IMessageChannel>(channel), IHostedService, IAsyncDisposable
{
    private Task _messageTask = null!;

    /// <summary>
    /// disposing functionality
    /// </summary>
    /// <returns></returns>
    public async ValueTask DisposeAsync()
    {
        await _channel.DisposeAsync();
        await _connection.DisposeAsync();
    }

    /// <summary>
    /// starting functionality
    /// </summary>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    public async Task StartAsync(CancellationToken cancellationToken)
    {
        await StartBroker();
        _messageTask = Task.Run(async () =>
        {
            await ComputeMessages();
        }, cancellationToken);
    }

    /// <summary>
    /// stopping functionality
    /// </summary>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
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