using DotNetEnv;
using RabbitMQ.Client;

namespace UserManagementSystem.Services.RabbitMQ;

public abstract class RabbitMQBase<Type> where Type : class, IRabbitChannel
{
    protected IConnection _connection = null!;
    protected IChannel _channel = null!;
    protected string URL { get; private set; } = null!;
    protected string Kind { get; private set; } = "direct";
    protected string Exchange { get; private set; } = "logger_service";
    protected string Queue { get; private set; } = "logs";
    protected string RoutingKey { get; private set; } = "logstore";
    protected readonly Type _messageChannel;

    public RabbitMQBase(Type channel)
    {
        _messageChannel = channel;

        var host = Env.GetString("RABBIT_HOST", "localhost");
        var port = Env.GetString("RABBIT_PORT", "5672");
        var user = Env.GetString("RABBIT_USER", "guest");
        var pass = Env.GetString("RABIT_PASSWORD", "guest");
        var vhost = Env.GetString("RABBIT_V_HOST", "my_vhost");

        URL = $"amqp://{user}:{pass}@{host}:{port}/{vhost}";
    }

    protected abstract Task StartBroker();
    protected abstract Task ComputeMessages();
}