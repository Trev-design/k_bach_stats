using System.Threading.Channels;

namespace UserManagementSystem.Services.RabbitMQ;

/// <summary>
/// 
/// </summary>
public class RabbitMessageChannel : IRabbitChannel, IMessageChannel, IMessagePipe
{
    private readonly Channel<byte[]> _channel = Channel.CreateUnbounded<byte[]>();

    public void Complete()
    {
        _channel.Writer.Complete();
    }

    public IAsyncEnumerable<byte[]> GetMessagePipe()
    {
        return _channel.Reader.ReadAllAsync();
    }

    public async Task SendMessageAsync(byte[] message)
    {
        await _channel.Writer.WriteAsync(message);
    }
}