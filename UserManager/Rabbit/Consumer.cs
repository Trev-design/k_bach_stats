using System.Text;
using System.Text.Json;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using UserManager.Redis.Data;
using UserManager.Redis.Models;

namespace UserManager.Rabbit;

public class Consumer
{
    private readonly IConnection _connection;
    private readonly IModel _channel;
    private string? _startSessionConsumerTag;
    private string? _stopSessionConsumerTag;
    private string? _addAccountConsumerTag;
    private string? _deleteAccountConsumerTag;
    private readonly UserHandler _userHandler;
    private readonly SessionHandler _sessionHandler;

    private int _messageCount = 0;

    public Consumer(IServiceScopeFactory scopeFactory)
    {
        _userHandler = new(scopeFactory);
        _sessionHandler = new(scopeFactory);

        var connectionFactory = new ConnectionFactory()
        {
            UserName = "IAmTheUser",
            Password = "ThisIsMyPassword",
            HostName = "localhost",
            Port = 5672,
            VirtualHost = "kbach"
        };

        _connection = connectionFactory.CreateConnection();

        _channel = _connection.CreateModel();
        _channel.ExchangeDeclare("session", "direct", true, false, null);

        _channel.ExchangeDeclare("account", "direct", true, false, null);

        _channel.QueueDeclare("start_user_session", true, false, false, null);
        _channel.QueueDeclare("stop_user_session", true, false, false, null);

        _channel.QueueDeclare("add_account", true, false, false, null);
        _channel.QueueDeclare("delete_account", true, false, false, null);

        _channel.QueueBind("start_user_session", "session", "send_session_credentials", null);
        _channel.QueueBind("stop_user_session", "session", "remove_user_session", null);

        _channel.QueueBind("add_account", "account", "add_account_request", null);
        _channel.QueueBind("delete_account", "account", "delete_account_request", null);
    }

    public async Task ConsumeSessionRequests(CancellationToken cancellationToken)
    {
        var consumer = new AsyncEventingBasicConsumer(_channel);
        consumer.Received += async (chan, args) => {
            await HandleStartSessionMessage(args);
        };

        var consumer2 = new AsyncEventingBasicConsumer(_channel);
        consumer2.Received += async (chan, args) => {
            await HandleStopSessionMessage(args);
        };

        var consumer3 = new AsyncEventingBasicConsumer(_channel);
        consumer3.Received += async (chan, args) => {
            await HandleAddAccountRequest(args);
        };

        var consumer4 = new AsyncEventingBasicConsumer(_channel);
        consumer4.Received += async (chan, args) => {
            await HandleDeleteAccountRequest(args);
        };

        _startSessionConsumerTag = _channel.BasicConsume("start_user_session", false, consumer);
        _stopSessionConsumerTag = _channel.BasicConsume("stop_user_session", false, consumer2);
        _addAccountConsumerTag = _channel.BasicConsume("add_account", false, consumer3);
        _deleteAccountConsumerTag = _channel.BasicConsume("delete_account", false, consumer4);

        await Task.Delay(Timeout.Infinite, cancellationToken);
    }

    public void StopConsumer()
    {
        Console.WriteLine("IAmOut");
        if (_startSessionConsumerTag != null)
        {
            _channel.BasicCancel(_startSessionConsumerTag);
        }
        
        if (_stopSessionConsumerTag != null)
        {
            _channel.BasicCancel(_stopSessionConsumerTag);
        }

        if (_addAccountConsumerTag != null)
        {
            _channel.BasicCancel(_addAccountConsumerTag);
        }

        if (_deleteAccountConsumerTag != null)
        {
            _channel.BasicCancel(_deleteAccountConsumerTag);
        }

        while (_messageCount > 0) {}

        _channel.Close();
        _connection.Close();
    }

    private async Task HandleStartSessionMessage(BasicDeliverEventArgs args)
    {
        Interlocked.Increment(ref _messageCount);

        try {
            await _sessionHandler.StartSession(args.Body);
            _channel.BasicAck(args.DeliveryTag, false);
        }
        catch (Exception ex) {
            _channel.BasicNack(args.DeliveryTag, false, false);
            Console.WriteLine(ex.Message);
        }
        finally {
            Interlocked.Decrement(ref _messageCount);
        }
    }

    private async Task HandleStopSessionMessage(BasicDeliverEventArgs args)
    {
        Interlocked.Increment(ref _messageCount);

        try {
            await _sessionHandler.StopSession(args.Body);
            _channel.BasicAck(args.DeliveryTag, false);
        }
        catch (Exception ex) {
            Console.WriteLine(ex.Message);
            _channel.BasicNack(args.DeliveryTag, false, false);
        }
        finally {
            Interlocked.Decrement(ref _messageCount);
        }
    }

    private async Task HandleAddAccountRequest(BasicDeliverEventArgs args)
    {
        Interlocked.Increment(ref _messageCount);

        try {
            await _userHandler.MakeUser(args.Body);
            _channel.BasicAck(args.DeliveryTag, false);
        }
        catch (Exception ex) {
            Console.WriteLine(ex.Message);
            _channel.BasicNack(args.DeliveryTag, false, false);
        }
        finally {
            Interlocked.Decrement(ref _messageCount);
        }
    }

    private async Task HandleDeleteAccountRequest(BasicDeliverEventArgs args)
    {
        Interlocked.Increment(ref _messageCount);

        try {
            await _userHandler.DeleteUser(args.Body);
            _channel.BasicAck(args.DeliveryTag, false);
        }
        catch (Exception ex) {
            Console.WriteLine(ex.Message);
            _channel.BasicNack(args.DeliveryTag, false, false);
        } 
        finally {
            Interlocked.Decrement(ref _messageCount);
        }
    }
}

