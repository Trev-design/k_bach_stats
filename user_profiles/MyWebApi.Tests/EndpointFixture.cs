
using System.Net;
using System.Security.Cryptography;
using Microsoft.AspNetCore;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Hosting.Server;
using Microsoft.AspNetCore.Hosting.Server.Features;
using Microsoft.AspNetCore.Server.Kestrel.Core;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using MyWebApi.Tests.Utils;
using UserManagementSystem.Services.Database;

namespace MyWebApi.Tests;

public class EndpointsFixture(DatabaseFixture dbFixture) : IAsyncLifetime
{
    private readonly DatabaseFixture _dbFixture = dbFixture;
    public HttpClient Client { get; private set; } = null!;
    public List<Guid> UserIDs { get; private set; } = null!;
    public List<Guid> ProfileIDs { get; private set; } = null!;
    public List<string> Entities { get; private set; } = null!;
    private IHost _host = null!;


    public async Task DisposeAsync()
    {
        await _host.StopAsync();
        _host.Dispose();
    }

    public async Task InitializeAsync()
    {
        _host = Host.CreateDefaultBuilder()
        .ConfigureWebHostDefaults(webBuilder =>
        {
            webBuilder.UseKestrel(options =>
            {
                options.Listen(IPAddress.Loopback, 0, options => options.Protocols = HttpProtocols.Http1AndHttp2);
            });

            webBuilder.ConfigureServices(services =>
            {
                services.AddControllers();
                services.AddDbContext<AppDBContext>(options =>
                {
                    options.UseMySql(_dbFixture.ConnectionString, ServerVersion.AutoDetect(_dbFixture.ConnectionString));
                });
            });

            webBuilder.Configure(app =>
            {
                app.UseRouting();
                app.UseEndpoints(endpoints =>
                {
                    endpoints.MapControllers();
                });
            });
        })
        .Build();

        await _host.StartAsync();

        var serverAddresses = _host.Services.GetRequiredService<IServer>().Features.Get<IServerAddressesFeature>();

        var address = serverAddresses!.Addresses.First();
        Client = new HttpClient { BaseAddress = new Uri(address) };
    }

    private async Task SetupDatabase()
    {
        Entities = [];

        for (int index = 0; index < 5; ++index)
        {
            var name = RandomString.GenerateRandomString(20);
            var email = RandomString.GenerateRandomEmail(15);
            var entity = Guid.NewGuid().ToString();
            Entities.Add(entity);

            await UserDBImpl.CreateUser(_dbFixture.Context, name, email, entity);
        }

        var users = await UserDBImpl.GetAllAsync(_dbFixture.Context);
        foreach (var user in users)
        {
            UserIDs.Add(user.Id);
            ProfileIDs.Add(user.UserProfile.Id);
        }
    }
}