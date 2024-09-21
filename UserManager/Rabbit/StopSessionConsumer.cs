
using RabbitMQ.Client.Events;

namespace UserManager.Rabbit;

public class StopSessionConsumer(string exchange, string queue, string route, RabbitConn conn, IServiceScopeFactory scopeFactory) : RabbitConsumer(exchange, queue, route, conn)
{
    private readonly SessionHandler _sessionHandler = new(scopeFactory);

    protected override async Task ComputeConsumeTask(BasicDeliverEventArgs args)
    {
        Interlocked.Increment(ref _messageCounter);

        try {
            await _sessionHandler.StopSession(args.Body);
            _channel.BasicAck(args.DeliveryTag, false);
        }
        catch (Exception ex) {
            _channel.BasicReject(args.DeliveryTag, false);
            Console.WriteLine(ex.Message);
        }
        finally {
            Interlocked.Decrement(ref _messageCounter);
        }
    }
}
