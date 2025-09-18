using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Testcontainers.MySql;
using UserManagementSystem.Models;
using UserManagementSystem.Services.Database;

namespace MyWebApi.Tests;

public class DatabaseFixture : IAsyncLifetime
{
    private MySqlContainer _container = null!;
    public AppDBContext Context { get; private set; } = null!;
    public string ConnectionString { get; private set; } = null!;

    public async Task InitializeAsync()
    {
        _container = new MySqlBuilder()
            .WithDatabase("test_db")
            .WithUsername("root")
            .WithPassword("testpass")
            .Build();
        // Container starten
        await _container.StartAsync();

        ConnectionString = _container.GetConnectionString();

        // DbContext erstellen
        var options = new DbContextOptionsBuilder<AppDBContext>()
            .UseMySql(_container.GetConnectionString(), ServerVersion.AutoDetect(_container.GetConnectionString()))
            .LogTo(Console.WriteLine, LogLevel.Warning)
            .Options;

        Context = new AppDBContext(options);

        // Migration nur einmal ausführen
        await Context.Database.MigrateAsync();
    }

    public async Task DisposeAsync()
    {
        await Context.DisposeAsync();
        await _container.DisposeAsync();
    }

    // Tabellen vor jedem Test leeren für Isolation
    public async Task ResetTablesAsync()
    {
        await Context.Database.ExecuteSqlRawAsync("SET FOREIGN_KEY_CHECKS=0;");
        await Context.Database.ExecuteSqlRawAsync("TRUNCATE TABLE Users;");
        await Context.Database.ExecuteSqlRawAsync("TRUNCATE TABLE Workspaces;");
        await Context.Database.ExecuteSqlRawAsync("TRUNCATE TABLE ChatRooms;");
        await Context.Database.ExecuteSqlRawAsync("TRUNCATE TABLE Profiles;");
        await Context.Database.ExecuteSqlRawAsync("TRUNCATE TABLE Contacts;");
        await Context.Database.ExecuteSqlRawAsync("SET FOREIGN_KEY_CHECKS=1;");
    }
}