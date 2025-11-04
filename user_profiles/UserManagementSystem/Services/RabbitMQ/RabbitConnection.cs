using DotNetEnv;
using RabbitMQ.Client;

namespace UserManagementSystem.Services.RabbitMQ;

public class RabbitConnection
{
    public string Host { get; } = Env.GetString("RABBIT_HOST", "localhost");
    public string Port { get; } = Env.GetString("RABBIT_PORT", "5672");
    public string User { get; } = Env.GetString("RABBIT_USER", "guest");
    public string Password { get; } = Env.GetString("RABbIT_PASSWORD", "guest");
    public string VirtualHost { get; } = Env.GetString("RABBIT_V_HOST", "my_vhost");

    public async Task<IConnection> MakeConnectionAsync()
    {
        var factory = new ConnectionFactory { Uri = new($"amqp://{User}:{Password}@{Host}:{Port}/{VirtualHost}") };
        return await factory.CreateConnectionAsync();
    }
}