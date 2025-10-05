using System.Threading.Channels;
using DotNetEnv;
using Microsoft.EntityFrameworkCore.Metadata.Conventions;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using UserManagementSystem.Services.RabbitMQ;

namespace MyWebApi.Tests.Mocks;

public abstract class ConsumerService(IMessageChannel channel) : RabbitMQBase(channel)
{
    protected IConnection _connection = null!;
    protected IChannel _channel = null!;
    protected readonly string url = "amqp://guest:guest@localhost:5672/my_vhost";
    protected string Tag { get; private set; } = null!;
    protected string Kind { get; private set; } = "direct";
    protected string Exchange { get; private set; } = "logger_service";
    protected string Queue { get; private set; } = "logs";
    protected string RoutingKey { get; private set; } = "logstore";

    protected abstract Task StartBroker();

    protected async override Task ComputeMessages()
    {
        var consumer = new AsyncEventingBasicConsumer(_channel);
        consumer.ReceivedAsync += async (model, args) =>
        {
            var body = args.Body.ToArray();
            await _messageChannel.SendMessageAsync(body);
        };

        Tag = await _channel.BasicConsumeAsync(Queue, true, consumer);
    } 
}