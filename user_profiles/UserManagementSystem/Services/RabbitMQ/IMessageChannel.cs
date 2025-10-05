namespace UserManagementSystem.Services.RabbitMQ;

public interface IMessageChannel
{
    public Task SendMessageAsync(byte[] message);
    public IAsyncEnumerable<byte[]> GetMessagePipe();
}