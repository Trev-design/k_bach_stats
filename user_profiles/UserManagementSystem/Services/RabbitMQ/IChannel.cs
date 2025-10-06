namespace UserManagementSystem.Services.RabbitMQ;

public interface IRabbitChannel
{
    public Task SendMessageAsync(byte[] message);
    public IAsyncEnumerable<byte[]> GetMessagePipe();
}

public interface IMessageChannel : IRabbitChannel { }

public interface IMessagePipe : IRabbitChannel {}