using Grpc.Net.Client;
using Microsoft.Extensions.DependencyInjection;
using UserManagementSystem.Grpc;
using UserManagementSystem.Services.Database;

namespace MyWebApi.Tests;

[Collection("Database collection")]
public class TestGRPC(GRPCTestFixture grpcFixture) : IClassFixture<GRPCTestFixture>
{
    private readonly GRPCTestFixture _grpcFixture = grpcFixture;

    [Fact]
    public async Task SavesUsers()
    {
        using var channel = GrpcChannel.ForAddress(_grpcFixture.Address);
        var client = new UserRegistryService.UserRegistryServiceClient(channel);

        using var call = client.UserPrimaryStream();
        var cts = new CancellationTokenSource(TimeSpan.FromSeconds(5));
        var responses = new List<RegistryResponse>();

        var readTask = Task.Run(async () =>
        {
            try
            {
                while (await call.ResponseStream.MoveNext(cts.Token))
                {
                    responses.Add(call.ResponseStream.Current);
                }
            }
            catch { }
        });

        await call.RequestStream.WriteAsync(new RegistryRequest { Name = "Alice", Email = "alice@test.com", Entity = "acme" });
        await call.RequestStream.WriteAsync(new RegistryRequest { Name = "Bob", Email = "bob@test.com", Entity = "bcme" });
        await call.RequestStream.WriteAsync(new RegistryRequest { Name = "Celine", Email = "celine@test.com", Entity = "ccme" });
        await call.RequestStream.WriteAsync(new RegistryRequest { Name = "Daniel", Email = "daniel@test.com", Entity = "dcme" });
        await call.RequestStream.CompleteAsync();
        await readTask;

        Assert.True(responses.All(response => response.Status == "ACCEPTED"));

        using var scope = _grpcFixture.Services.CreateScope();
        var db = scope.ServiceProvider.GetRequiredService<AppDBContext>();
        var users = await UserDBImpl.GetAllAsync(db);

        Console.WriteLine(db);

        Assert.Contains(users, user => user.UserProfile.UserContact.Name == "Alice");
        Assert.Contains(users, user => user.UserProfile.UserContact.Name == "Bob");
        Assert.Contains(users, user => user.UserProfile.UserContact.Name == "Celine");
        Assert.Contains(users, user => user.UserProfile.UserContact.Name == "Daniel");
    }

    [Fact]
    public async Task SaveUsersFailedDoubleEntitySimultaniously()
    {
        using var channel = GrpcChannel.ForAddress(_grpcFixture.Address);
        var client = new UserRegistryService.UserRegistryServiceClient(channel);

        using var call = client.UserPrimaryStream();
        var cts = new CancellationTokenSource(TimeSpan.FromSeconds(5));
        var responses = new List<RegistryResponse>();

        var readTask = Task.Run(async () =>
        {
            try
            {
                while (await call.ResponseStream.MoveNext(cts.Token))
                {
                    responses.Add(call.ResponseStream.Current);
                }
            }
            catch { }
        });

        await call.RequestStream.WriteAsync(new RegistryRequest { Name = "Garreth", Email = "gary@test.com", Entity = "ecme" });
        await call.RequestStream.WriteAsync(new RegistryRequest { Name = "Herb", Email = "herb@test.com", Entity = "ecme" });
        await call.RequestStream.CompleteAsync();
        await readTask;

        Assert.Contains(responses, response => response.Status == "DUPLICATE");
    }


    [Fact]
    public async Task SaveUserFailedDoubleEmailSimultaniously()
    {
        using var channel = GrpcChannel.ForAddress(_grpcFixture.Address);
        var client = new UserRegistryService.UserRegistryServiceClient(channel);

        using var call = client.UserPrimaryStream();
        var cts = new CancellationTokenSource(TimeSpan.FromSeconds(5));
        var responses = new List<RegistryResponse>();

        var readTask = Task.Run(async () =>
        {
            try
            {
                while (await call.ResponseStream.MoveNext(cts.Token))
                {
                    responses.Add(call.ResponseStream.Current);
                }
            }
            catch { }
        });

        await call.RequestStream.WriteAsync(new RegistryRequest { Name = "Irene", Email = "irene@test.com", Entity = "hcme" });
        await call.RequestStream.WriteAsync(new RegistryRequest { Name = "Jasmine", Email = "irene@test.com", Entity = "icme" });
        await call.RequestStream.CompleteAsync();
        await readTask;

        Assert.Contains(responses, response => response.Status == "DUPLICATE");
    }
}