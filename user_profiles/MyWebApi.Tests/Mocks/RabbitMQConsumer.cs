using Microsoft.Extensions.Hosting;
using RabbitMQ.Client;
using UserManagementSystem.Services.RabbitMQ;

namespace MyWebApi.Tests.Mocks;

public sealed class RabbitMQConsumer(IMessageChannel channel) : ConsumerService(channel), IHostedService, IAsyncDisposable
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
        await _channel.CloseAsync(cancellationToken: cancellationToken);
        await _connection.CloseAsync(cancellationToken: cancellationToken);
    }

    protected override async Task StartBroker()
    {
        var factory = new ConnectionFactory { Uri = new(url) };
        _connection = await factory.CreateConnectionAsync();
        _channel = await _connection.CreateChannelAsync();

        await _channel.ExchangeDeclareAsync(Exchange, Kind, true, false, null, false);
        await _channel.QueueDeclareAsync(Queue, true, false, false, null, false);
        await _channel.QueueBindAsync(Queue, Exchange, RoutingKey, null, false);;
    }
}