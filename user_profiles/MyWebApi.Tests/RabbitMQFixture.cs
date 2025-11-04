using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using MyWebApi.Tests.Mocks;
using Testcontainers.RabbitMq;
using UserManagementSystem.Services.RabbitMQ;

namespace MyWebApi.Tests;

public class RabbitMQFixture : IAsyncLifetime
{
    private IHost _host = null!;
    private RabbitMqContainer _container = null!;
    public IServiceProvider Services { get; private set; } = null!;

    public async Task DisposeAsync()
    {
        await _host.StopAsync();
        await _container.StopAsync();
        await _container.DisposeAsync();
    }

    public async Task InitializeAsync()
    {
        _container = new RabbitMqBuilder()
        .WithImage("rabbitmq:4.1.1-management")
        .WithHostname("localhost")
        .WithPortBinding(5672)
        .WithUsername("guest")
        .WithPassword("guest")
        .WithEnvironment("RABBITMQ_DEFAULT_VHOST", "my_vhost")
        .Build();

        await _container.StartAsync();

        _host = Host.CreateDefaultBuilder()
        .ConfigureServices(services =>
        {
            services.AddSingleton<ILogMessageChannel, RabbitMessageChannel>();
            services.AddSingleton<ILogMessagePipe, RabbitMessageChannel>();
            services.AddSingleton<ISearchMessageChannel, RabbitMessageChannel>();
            services.AddSingleton<ISearchMessagePipe, RabbitMessageChannel>();
            services.AddHostedService<RabbitMQLoggingService>();
            services.AddHostedService<RabbitMQLogConsumer>();
            services.AddHostedService<RabbitMQSearchEngineConsumer>();
        })
        .Build();

        await _host.StartAsync();

        Services = _host.Services;
    }
}