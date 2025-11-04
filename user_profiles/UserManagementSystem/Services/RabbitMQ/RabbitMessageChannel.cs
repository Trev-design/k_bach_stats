using System.Threading.Channels;
 
namespace UserManagementSystem.Services.RabbitMQ;

/// <summary>
/// channel for local traffic
/// </summary>
public class RabbitMessageChannel
: IRabbitChannel
, ILogMessageChannel
, ISearchMessageChannel
, ILogMessagePipe
, ISearchMessagePipe
{
    private readonly Channel<byte[]> _channel = Channel.CreateUnbounded<byte[]>();

    /// <summary>
    /// closes the channel
    /// </summary>
    public void Complete()
    {
        _channel.Writer.Complete();
    }

    /// <summary>
    /// gives you an enumerable to iterate over the channel
    /// </summary>
    /// <returns></returns>
    public IAsyncEnumerable<byte[]> GetMessagePipe()
    {
        return _channel.Reader.ReadAllAsync();
    }

    /// <summary>
    /// here you can send your messages
    /// </summary>
    /// <param name="message"></param>
    /// <returns></returns>
    public async Task SendMessageAsync(byte[] message)
    {
        await _channel.Writer.WriteAsync(message);
    }
}

public class RabbitLogMessageChannel : RabbitMessageChannel { }
public class RabbitSearchMessageChannel : RabbitMessageChannel { }
public class RabbitLogMessagePipe : RabbitMessageChannel { }
public class RabbitSearchMessagePipe : RabbitMessageChannel { }