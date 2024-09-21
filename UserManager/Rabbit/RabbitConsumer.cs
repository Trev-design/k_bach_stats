using RabbitMQ.Client;
using RabbitMQ.Client.Events;

namespace UserManager.Rabbit;

public abstract class RabbitConsumer
{
    protected readonly IModel _channel;
    private readonly string _queue;
    private string? _consumerTag;
    protected int _messageCounter = 0;

    public RabbitConsumer(string exchange, string queue, string route, RabbitConn conn)
    {
        _queue = queue;
        _channel = conn.OpenChannel();
        _channel.ExchangeDeclare(exchange, "direct", true, false);
        _channel.QueueDeclare(queue, true, false, false);
        _channel.QueueBind(queue, exchange, route, null);
    }

    public Task Consume(CancellationToken cancellationToken)
    {
        var consumer = new AsyncEventingBasicConsumer(_channel);
        consumer.Received += async (chan, args) => {
            await ComputeConsumeTask(args);
        };

        _consumerTag = _channel.BasicConsume(_queue, false, "", true, false, null, consumer);

        return Task.CompletedTask;
    }

    public Task StopConsumer()
    {
        if (_consumerTag != null)
        {
            _channel.BasicCancel(_consumerTag);
        }

        while (_messageCounter > 0) {}

        _channel.Close();

        return Task.CompletedTask;
    }

    protected abstract Task ComputeConsumeTask(BasicDeliverEventArgs args);
}

