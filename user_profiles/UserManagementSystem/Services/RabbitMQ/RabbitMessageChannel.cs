using System.Threading.Channels;

namespace UserManagementSystem.Services.RabbitMQ;

public class RabbitMessageChannel : IMessageChannel
{
    private readonly Channel<byte[]> _channel = Channel.CreateUnbounded<byte[]>();

    public IAsyncEnumerable<byte[]> GetMessagePipe()
    {
        return _channel.Reader.ReadAllAsync();
    }

    public async Task SendMessageAsync(byte[] message)
    {
       await _channel.Writer.WriteAsync(message);
    }
}