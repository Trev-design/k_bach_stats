using Microsoft.EntityFrameworkCore;
using Testcontainers.MySql;
using UserManagementSystem.Models;
using UserManagementSystem.Services;

namespace MyWebApi.Tests;

public class DatabaseFixture : IAsyncLifetime
{
    private readonly MySqlContainer _container;
    public AppDBContext Context { get; private set; } = null!;

    public DatabaseFixture()
    {
        _container = new MySqlBuilder()
            .WithDatabase("test_db")
            .WithUsername("root")
            .WithPassword("testpass")
            .Build();
    }

    public async Task InitializeAsync()
    {
        // Container starten
        await _container.StartAsync();

        // DbContext erstellen
        var options = new DbContextOptionsBuilder<AppDBContext>()
            .UseMySql(_container.GetConnectionString(), ServerVersion.AutoDetect(_container.GetConnectionString()))
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