using RabbitMQ.Client;

namespace UserManagementSystem.Services.RabbitMQ;

public class RabbitChannel
{
    public required string Exchange { get; init; }
    public required string Queue { get; init; }
    public required string Type { get; init; }
    public required string Key { get; init; }

    public async Task<IChannel> MakeChannelAsync(IConnection connection)
    {
        var channel = await connection.CreateChannelAsync();
        await channel.ExchangeDeclareAsync(Exchange, Type, true, false, null, false);
        await channel.QueueDeclareAsync(Queue, true, false, false, null, false);
        await channel.QueueBindAsync(Queue, Exchange, Key, null, false);

        return channel;
    }
} 