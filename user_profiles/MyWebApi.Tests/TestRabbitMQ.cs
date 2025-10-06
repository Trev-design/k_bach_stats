using System.Text;
using Microsoft.Extensions.DependencyInjection;
using UserManagementSystem.Services.RabbitMQ;

namespace MyWebApi.Tests;

public class TestRabbitMQ(RabbitMQFixture fixture) : IClassFixture<RabbitMQFixture>
{
    private readonly RabbitMQFixture _fixture = fixture;

    [Fact]
    public async Task TestRabbitMQInfrastructure()
    {
        var channel = _fixture.Services.GetRequiredService<IMessageChannel>();
        var pipe = _fixture.Services.GetRequiredService<IMessagePipe>();
        List<string> responses = [];
        List<byte[]> requests = [Encoding.UTF8.GetBytes("hello"), Encoding.UTF8.GetBytes("miss"), Encoding.UTF8.GetBytes("jackson")];

        var task = Task.Run(async () =>
        {
            await foreach (var message in pipe.GetMessagePipe())
            {
                responses.Add(Encoding.UTF8.GetString(message));
            }
        });

        foreach (var request in requests)
        {
            await channel.SendMessageAsync(request);
        }

        await Task.Delay(100);
        pipe.Complete();
        await task;

        Assert.Equal(3, responses.Count);
        Assert.Contains(responses, response => response == "hello");
        Assert.Contains(responses, response => response == "miss");
        Assert.Contains(responses, response => response == "jackson");
    }
}