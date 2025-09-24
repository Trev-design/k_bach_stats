using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.AspNetCore.TestHost;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using MyWebApi.Tests.Utils;
using UserManagementSystem.Services.Database;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Builder;
using UserManagementSystem.Controllers;
using Microsoft.Extensions.Hosting;

namespace MyWebApi.Tests;

public sealed class EndpointsFixture(DatabaseFixture dbFixture) : IAsyncLifetime
{
    private readonly DatabaseFixture _dbFixture = dbFixture;
    public HttpClient Client { get; private set; } = null!;
    public List<Guid> UserIDs { get; private set; } = [];
    public List<Guid> DeleteUserIDs { get; private set; } = [];
    public Dictionary<Guid, Guid> Workspaces { get; set; } = [];
    public Dictionary<Guid, Guid> ChatRooms { get; set; } = [];
    public List<Guid> ProfileIDs { get; private set; } = [];
    public List<string> Entities { get; private set; } = [];
    private IHost _host = null!;

    public async Task DisposeAsync()
    {
        if (_host != null)
        {
            await _host.StopAsync();
            _host.Dispose();
        }
    }

    public async Task InitializeAsync()
    {
        var builder = WebApplication.CreateBuilder();

        builder.Services.AddDbContext<AppDBContext>(options =>
        {
            options.UseMySql(_dbFixture.ConnectionString, ServerVersion.AutoDetect(_dbFixture.ConnectionString));
        });

        builder.Services.AddControllers()
            .AddApplicationPart(typeof(UserController).Assembly)
            .AddApplicationPart(typeof(ProfileController).Assembly);

        builder.WebHost.UseTestServer();

        var app = builder.Build();

        app.MapControllers();

        await app.StartAsync();

        Client = app.GetTestClient();

        _host = app;
    }

    private async Task SetupDatabase()
    {
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