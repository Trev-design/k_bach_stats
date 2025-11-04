using DotNetEnv;
using RabbitMQ.Client;
using UserManagementSystem.Services.RabbitMQ;

namespace MyWebApi.Tests.Mocks;

/// <summary>
/// the baseclass of your rabbitmq infrastructure
/// </summary>
/// <typeparam name="Type"></typeparam>
public abstract class RabbitMQBase<Type> where Type : class, IRabbitChannel
{
    protected IConnection _connection = null!;
    protected IChannel _channel = null!;
    protected string URL { get; private set; } = null!;
    protected string Kind { get; private set; } = "direct";
    protected string Exchange { get; private set; } = null!;
    protected string Queue { get; private set; } = null!;
    protected string RoutingKey { get; private set; } = null!;
    protected readonly Type _messageChannel;

    public RabbitMQBase(Type channel, string exchange, string queue, string key) 
    {
        Exchange = exchange;
        Queue = queue;
        RoutingKey = key;
        _messageChannel = channel;

        var host = Env.GetString("RABBIT_HOST", "localhost");
        var port = Env.GetString("RABBIT_PORT", "5672");
        var user = Env.GetString("RABBIT_USER", "guest");
        var pass = Env.GetString("RABbIT_PASSWORD", "guest");
        var vhost = Env.GetString("RABBIT_V_HOST", "my_vhost");

        URL = $"amqp://{user}:{pass}@{host}:{port}/{vhost}";
    }

    protected async Task StartBroker()
    {
        var factory = new ConnectionFactory { Uri = new(URL) };
        _connection = await factory.CreateConnectionAsync();
        _channel = await _connection.CreateChannelAsync();

        await _channel.ExchangeDeclareAsync(Exchange, Kind, true, false, null, false);
        await _channel.QueueDeclareAsync(Queue, true, false, false, null, false);
        await _channel.QueueBindAsync(Queue, Exchange, RoutingKey, null, false);
    }
    
    protected abstract Task ComputeMessages();
}