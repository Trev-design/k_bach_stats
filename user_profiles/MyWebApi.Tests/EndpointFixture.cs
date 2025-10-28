using Microsoft.AspNetCore.TestHost;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using MyWebApi.Tests.Utils;
using UserManagementSystem.Services.Database;
using Microsoft.AspNetCore.Builder;
using UserManagementSystem.Controllers;
using Microsoft.Extensions.Hosting;
using UserManagementSystem.Models;
using System.Text.Json.Serialization;

namespace MyWebApi.Tests;

public sealed class EndpointsFixture(DatabaseFixture dbFixture) : IAsyncLifetime
{
    private readonly DatabaseFixture _dbFixture = dbFixture;
    public HttpClient Client { get; private set; } = null!;
    public List<Guid> UserIDs { get; private set; } = [];
    public List<Guid> UserIDsWithWorkspaces { get; private set; } = [];
    public List<Guid> UserIDsWithChatrooms { get; private set; } = [];
    public List<Guid> DeleteWorcspaceIDs { get; private set; } = [];
    public List<Guid> DeleteUserIDs { get; private set; } = [];
    public List<Guid> DeleteChatIDs { get; private set; } = [];
    public List<string> Entities { get; private set; } = [];
    public List<string> DeleteEntities { get; private set; } = [];
    public List<string> EntitiesWithWorkspaces { get; private set; } = [];
    public List<string> DeleteEntitiesWithWorkspaces { get; private set; } = [];
    public List<string> EntitiesWithChatRooms { get; private set; } = [];
    public List<string> DeleteEntitiesWithChatRooms { get; private set; } = [];

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

        builder.Services.AddControllers().AddJsonOptions(opt =>
        {
            opt.JsonSerializerOptions.ReferenceHandler = ReferenceHandler.IgnoreCycles;
        })
            .AddApplicationPart(typeof(UserController).Assembly)
            .AddApplicationPart(typeof(ProfileController).Assembly);

        builder.WebHost.UseTestServer();

        var app = builder.Build();

        app.MapControllers();

        await app.StartAsync();

        Client = app.GetTestClient();

        _host = app;

        await SetupDatabase();
    }

    private async Task SetupDatabase()
    {
        await GetRawUserIds(UserIDs, _dbFixture.Context);
        await GetRawUserIds(DeleteUserIDs, _dbFixture.Context);
        await GetUserWithWorkspaces(UserIDsWithWorkspaces, _dbFixture.Context);
        await GetUserWithWorkspaces(DeleteWorcspaceIDs, _dbFixture.Context);
        await GetUserIDsWithChats(UserIDsWithChatrooms, _dbFixture.Context);
        await GetUserIDsWithChats(DeleteChatIDs, _dbFixture.Context);
        await GetRawEntities(Entities, _dbFixture.Context);
        await GetRawEntities(DeleteEntities, _dbFixture.Context);
        await GetEntitiesWithWorkspaces(EntitiesWithWorkspaces, _dbFixture.Context);
        await GetEntitiesWithWorkspaces(DeleteEntitiesWithWorkspaces, _dbFixture.Context);
        await GetEntitiesWithChats(EntitiesWithChatRooms, _dbFixture.Context);
        await GetEntitiesWithChats(DeleteEntitiesWithChatRooms, _dbFixture.Context);
    }

    private static async Task GetUserWithWorkspaces(List<Guid> list, AppDBContext context)
    {
        for (var index = 0; index < 5; ++index)
        {
            var user = new User { Entity = RandomString.GenerateRandomString(36) };
            await context.Users.AddAsync(user);
            var profile = new Profile { UserId = user.Id };
            await context.Profiles.AddAsync(profile);
            var contact = new Contact
            {
                ProfileId = profile.Id,
                Email = RandomString.GenerateRandomEmail(15),
                Name = RandomString.GenerateRandomString(12)
            };
            await context.Contacts.AddAsync(contact);
            var workspace = new Workspace { UserId = user.Id, Name = RandomString.GenerateRandomString(22) };
            await context.Workspaces.AddAsync(workspace);
            await context.SaveChangesAsync();

            list.Add(user.Id);
        }
    }

    private static async Task GetRawUserIds(List<Guid> list, AppDBContext context)
    {
        for (var index = 0; index < 5; ++index)
        {
            var user = new User { Entity = RandomString.GenerateRandomString(36) };
            await context.Users.AddAsync(user);
            var profile = new Profile { UserId = user.Id };
            await context.Profiles.AddAsync(profile);
            var contact = new Contact
            {
                ProfileId = profile.Id,
                Email = RandomString.GenerateRandomEmail(15),
                Name = RandomString.GenerateRandomString(12)
            };
            await context.Contacts.AddAsync(contact);

            list.Add(user.Id);
        }
    }

    private static async Task GetUserIDsWithChats(List<Guid> list, AppDBContext context)
    {
        for (var index = 0; index < 5; ++index)
        {
            var user = new User { Entity = RandomString.GenerateRandomString(36) };
            await context.Users.AddAsync(user);
            var profile = new Profile { UserId = user.Id };
            await context.Profiles.AddAsync(profile);
            var contact = new Contact
            {
                ProfileId = profile.Id,
                Email = RandomString.GenerateRandomEmail(15),
                Name = RandomString.GenerateRandomString(12)
            };
            await context.Contacts.AddAsync(contact);
            var workspace = new Workspace { UserId = user.Id, Name = RandomString.GenerateRandomString(22) };
            await context.Workspaces.AddAsync(workspace);
            var chatRoom = new ChatRoom
            {
                WorkspaceId = workspace.Id,
                Reference = RandomString.GenerateRandomString(36),
                Topic = RandomString.GenerateRandomString(50)
            };
            await context.ChatRooms.AddAsync(chatRoom);
            await context.SaveChangesAsync();

            list.Add(user.Id);
        }
    }

    private static async Task GetEntitiesWithWorkspaces(List<string> list, AppDBContext context)
    {
        for (var index = 0; index < 5; ++index)
        {
            var user = new User { Entity = RandomString.GenerateRandomString(36) };
            await context.Users.AddAsync(user);
            var profile = new Profile { UserId = user.Id };
            await context.Profiles.AddAsync(profile);
            var contact = new Contact
            {
                ProfileId = profile.Id,
                Email = RandomString.GenerateRandomEmail(15),
                Name = RandomString.GenerateRandomString(12)
            };
            await context.Contacts.AddAsync(contact);
            var workspace = new Workspace { UserId = user.Id, Name = RandomString.GenerateRandomString(22) };
            await context.Workspaces.AddAsync(workspace);
            await context.SaveChangesAsync();

            list.Add(user.Entity);
        }
    }

    private static async Task GetRawEntities(List<string> list, AppDBContext context)
    {
        for (var index = 0; index < 5; ++index)
        {
            var user = new User { Entity = RandomString.GenerateRandomString(36) };
            await context.Users.AddAsync(user);
            var profile = new Profile { UserId = user.Id };
            await context.Profiles.AddAsync(profile);
            var contact = new Contact
            {
                ProfileId = profile.Id,
                Email = RandomString.GenerateRandomEmail(15),
                Name = RandomString.GenerateRandomString(12)
            };
            await context.Contacts.AddAsync(contact);

            list.Add(user.Entity);
        }
    }

    private static async Task GetEntitiesWithChats(List<string> list, AppDBContext context)
    {
        for (var index = 0; index < 5; ++index)
        {
            var user = new User { Entity = RandomString.GenerateRandomString(36) };
            await context.Users.AddAsync(user);
            var profile = new Profile { UserId = user.Id };
            await context.Profiles.AddAsync(profile);
            var contact = new Contact
            {
                ProfileId = profile.Id,
                Email = RandomString.GenerateRandomEmail(15),
                Name = RandomString.GenerateRandomString(12)
            };
            await context.Contacts.AddAsync(contact);
            var workspace = new Workspace { UserId = user.Id, Name = RandomString.GenerateRandomString(22) };
            await context.Workspaces.AddAsync(workspace);
            var chatRoom = new ChatRoom 
            {
                WorkspaceId = workspace.Id,
                Reference = RandomString.GenerateRandomString(36),
                Topic = RandomString.GenerateRandomString(50)
            };
            await context.ChatRooms.AddAsync(chatRoom);
            await context.SaveChangesAsync();

            list.Add(user.Entity);
        }
    }
}