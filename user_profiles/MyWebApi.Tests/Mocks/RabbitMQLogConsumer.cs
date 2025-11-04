using Microsoft.Extensions.Hosting;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using UserManagementSystem.Services.RabbitMQ;

namespace MyWebApi.Tests.Mocks;

public sealed class RabbitMQLogConsumer(ILogMessagePipe channel) 
: RabbitMQBase<ILogMessagePipe>(channel, "logger_service", "logs", "logstore")
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
        await _channel.CloseAsync(cancellationToken: cancellationToken);
        await _connection.CloseAsync(cancellationToken: cancellationToken);
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