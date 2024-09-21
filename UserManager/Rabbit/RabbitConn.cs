using RabbitMQ.Client;

namespace UserManager.Rabbit;

public sealed class RabbitConn : IDisposable
{
    private readonly IConnection _connection;

    public RabbitConn()
    {
        var connectionFactory = new ConnectionFactory()
        {
            UserName = "IAmTheUser",
            Password = "ThisIsMyPassword",
            HostName = "localhost",
            Port = 5672,
            VirtualHost = "kbach",
            DispatchConsumersAsync = true,
        };

        _connection = connectionFactory.CreateConnection();
    }

    public IModel OpenChannel() => _connection.CreateModel();

    public void Dispose()
    {
        _connection.Close();
    }
}