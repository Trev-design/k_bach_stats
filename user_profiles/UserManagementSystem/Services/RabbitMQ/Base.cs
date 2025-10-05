namespace UserManagementSystem.Services.RabbitMQ;

public abstract class RabbitMQBase(IMessageChannel channel)
{
    protected readonly IMessageChannel _messageChannel = channel;

    protected abstract Task ComputeMessages();
}