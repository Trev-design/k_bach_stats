using System.Text.Json;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using UserManager.Redis.Data;
using UserManager.Redis.Models;

namespace UserManager.Rabbit;

public class Consumer
{
    private readonly ISessionRepo _redis;
    private readonly IConnection _connection;
    private readonly IModel _channel;
    private string? _startSessionConsumerTag;
    private string? _stopSessionConsumerTag;

    public Consumer(ISessionRepo redis)
    {
        _redis = redis;

        var connectionFactory = new ConnectionFactory()
        {
            UserName = "IAmTheUser",
            Password = "ThisIsMyPassword",
            HostName = "localhost",
            Port = 5672,
            VirtualHost = "kbach"
        };

        _connection = connectionFactory.CreateConnection();

        _channel = _connection.CreateModel();
    }

    public async Task ConsumeSessionRequests(CancellationToken cancellationToken)
    {
        await Task.Run(() => {
            var consumer = new AsyncEventingBasicConsumer(_channel);
            consumer.Received += async (chan, args) => {
                await HandleStartSessionMessage(args);
            };

            var consumer2 = new AsyncEventingBasicConsumer(_channel);
            consumer2.Received += async (chan, args) => {
                await HandleStopSessionMessage(args);
            };

            _startSessionConsumerTag = _channel.BasicConsume("start_user_session", false, consumer);
            _stopSessionConsumerTag = _channel.BasicConsume("stop_user_session", false, consumer2);
        }, cancellationToken);

    }

    public void StopConsumer()
    {
        _channel.BasicCancel(_startSessionConsumerTag);
        _channel.BasicCancel(_stopSessionConsumerTag);
        _channel.Close();
        _connection.Close();
    }

    private async Task HandleStartSessionMessage(BasicDeliverEventArgs args)
    {
        var body = args.Body.ToString();
        Session? session = JsonSerializer.Deserialize<Session>(body);
        if (session == null) return;
        await _redis.CreateSessionAsync(session);
        _channel.BasicAck(args.DeliveryTag, false);
    }

    private async Task HandleStopSessionMessage(BasicDeliverEventArgs args)
    {
        var body = args.Body.ToString();
        Session? session = JsonSerializer.Deserialize<Session>(body);

        if (session == null) 
        {
            _channel.BasicNack(args.DeliveryTag, false, false);
            return;
        }

        await _redis.GetSession(session.ID);

        if (session == null) 
        {
            _channel.BasicNack(args.DeliveryTag, false, false);
            return;
        }

        await _redis.DeleteSessionAsync(session.ID);
        _channel.BasicAck(args.DeliveryTag, false);
    }
}

