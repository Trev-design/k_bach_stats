using System.Threading.Channels;
using RabbitMQ.Client;

namespace UserManagementSystem.Services.RabbitMQ;

/// <summary>
/// the rabbitmq implementation
/// </summary>
/// <param name="channel"></param>
public sealed class RabbitMQLoggingService(
    ILogMessageChannel logMessageChannel,
    ISearchMessageChannel searchMessageChannel
) : IHostedService, IAsyncDisposable
{
    private Task _logMessageTask = null!;
    private Task _searchEngineMessageTask = null!;
    private readonly ILogMessageChannel _logMessageChannel = logMessageChannel;
    private readonly ISearchMessageChannel _searchMessageChannel = searchMessageChannel;
    private IConnection _connection = null!;
    private IChannel _logChannel = null!;
    private IChannel _searchEngineChannel = null!;
    private readonly string _kind  = "direct";
    private readonly string _logExchange = "logger_service";
    private readonly string _searchEngineExchange = "search_engine_service";
    private readonly string _logQueue = "logs";
    private readonly string _searchEngineQueue = "search_engine_entities";
    private readonly string _logKey = "log_store";
    private readonly string _searchEngineKey = "search_engine";

    /// <summary>
    /// starting functionality
    /// </summary>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    public async Task StartAsync(CancellationToken cancellationToken)
    {
        await StartBroker();
        _searchEngineMessageTask = Task.Run(async () => await ComputeSearchEngineMessages(), cancellationToken);
        _logMessageTask = Task.Run(async () => ComputeLogMessages(), cancellationToken);
    }

    /// <summary>
    /// stopping functionality
    /// </summary>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    public async Task StopAsync(CancellationToken cancellationToken)
    {
        _logMessageChannel.Complete();
        _searchMessageChannel.Complete();
        await _searchEngineMessageTask;
        await _logMessageTask;
        await _logChannel.CloseAsync(cancellationToken);
        await _searchEngineChannel.CloseAsync(cancellationToken);
        await _connection.CloseAsync(cancellationToken);
    }

    /// <summary>
    /// dispose functionality
    /// </summary>
    /// <returns></returns>
    public async ValueTask DisposeAsync()
    {
        await _logChannel.DisposeAsync();
        await _searchEngineChannel.DisposeAsync();
        await _connection.DisposeAsync();
    }

    private async Task ComputeLogMessages()
    {
        await foreach (var message in _logMessageChannel.GetMessagePipe())
        {
            await _logChannel.BasicPublishAsync(_logExchange, _logKey, message);
        }
    }

    private async Task ComputeSearchEngineMessages()
    {
        await foreach (var message in _searchMessageChannel.GetMessagePipe())
        {
            await _searchEngineChannel.BasicPublishAsync(_searchEngineExchange, _searchEngineKey, message);
        }
    }

    private async Task StartBroker()
    {
        _connection = await new RabbitConnection { }.MakeConnectionAsync();

        _logChannel = await new RabbitChannel
        {
            Exchange = _logExchange,
            Queue = _logQueue,
            Type = _kind,
            Key = _logKey
        }.MakeChannelAsync(_connection);

        _searchEngineChannel = await new RabbitChannel
        {
            Exchange = _searchEngineExchange,
            Queue = _searchEngineQueue,
            Type = _kind,
            Key = _searchEngineKey
        }.MakeChannelAsync(_connection); 
    }
}