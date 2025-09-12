
using System.Net;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Hosting.Server;
using Microsoft.AspNetCore.Hosting.Server.Features;
using Microsoft.AspNetCore.Server.Kestrel.Core;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using UserManagementSystem.Services.Database;
using UserManagementSystem.Services.GRPC;

namespace MyWebApi.Tests;

public class GRPCTestFixture(DatabaseFixture dbFixture) : IAsyncLifetime
{
    private readonly string _dbConnstring = dbFixture.ConnectionString;
    private IHost _host = null!;
    public IServiceProvider Services => _host.Services;
    public string Address { get; private set; } = null!;

    public async Task DisposeAsync()
    {
        await _host.StopAsync();
    }

    public async Task InitializeAsync()
    {
        _host = Host.CreateDefaultBuilder().ConfigureWebHostDefaults(webBuilder =>
        {
            webBuilder.UseKestrel(options =>
            {
                options.Listen(IPAddress.Loopback, 0, options => options.Protocols = HttpProtocols.Http2);
            });

            webBuilder.ConfigureServices(services =>
            {
                services.AddGrpc();
                services.AddDbContext<AppDBContext>(context => context.UseMySql(_dbConnstring, ServerVersion.AutoDetect(_dbConnstring)));
                services.AddScoped<UserRegistryServiceImpl>();
            });

            webBuilder.Configure(app =>
            {
                app.UseRouting();
                app.UseEndpoints(endpoints => endpoints.MapGrpcService<UserRegistryServiceImpl>());
            });
        }).Build();

        await _host.StartAsync();

        var serverAddresses = _host.Services.GetRequiredService<IServer>().Features.Get<IServerAddressesFeature>();
        Address = serverAddresses!.Addresses.First();
    }
}