using RabbitMQ.Client.Events;

namespace UserManager.Rabbit;

public class AddUserConsumer(string exchange, string queue, string route, RabbitConn conn, IServiceScopeFactory scopeFactory) : RabbitConsumer(exchange, queue, route, conn)
{
    private readonly UserHandler _userHandler = new(scopeFactory);

    protected override async Task ComputeConsumeTask(BasicDeliverEventArgs args)
    {
        Interlocked.Increment(ref _messageCounter);
        Console.WriteLine("hey dude i'm here to add a user and its super duper cool");

        try {
            await _userHandler.MakeUser(args.Body);
            _channel.BasicAck(args.DeliveryTag, false);
            Console.WriteLine("you are in man");
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
