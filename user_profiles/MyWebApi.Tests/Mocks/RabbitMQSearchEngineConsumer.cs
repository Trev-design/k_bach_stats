using Microsoft.Extensions.Hosting;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using UserManagementSystem.Services.RabbitMQ;

namespace MyWebApi.Tests.Mocks;

public sealed class RabbitMQSearchEngineConsumer(ISearchMessageChannel channel)
: RabbitMQBase<ISearchMessageChannel>(channel, "search_engine_service", "search_engine_entities", "search_engine")
, IHostedService 
, IAsyncDisposable
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

    protected async override Task ComputeMessages()
    {
        var consumer = new AsyncEventingBasicConsumer(_channel);
        consumer.ReceivedAsync += async (model, args) =>
        {
            var body = args.Body.ToArray();
            await _messageChannel.SendMessageAsync(body);
        };

        await _channel.BasicConsumeAsync(Queue, true, consumer);
    }
}