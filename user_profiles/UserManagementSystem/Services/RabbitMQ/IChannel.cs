namespace UserManagementSystem.Services.RabbitMQ;

/// <summary>
/// interface for local traffic in your rabbitmq infrastructure
/// </summary>
public interface IRabbitChannel
{
    public Task SendMessageAsync(byte[] message);
    public IAsyncEnumerable<byte[]> GetMessagePipe();
    public void Complete();
}

public interface ILogMessageChannel : IRabbitChannel { }

public interface ILogMessagePipe : IRabbitChannel { }

public interface ISearchMessageChannel : IRabbitChannel { }

public interface ISearchMessagePipe : IRabbitChannel {}