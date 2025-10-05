namespace MyWebApi.Tests;

public class TestRabbitMQ(RabbitMQFixture fixture) : IClassFixture<RabbitMQFixture>
{
    private readonly RabbitMQFixture _fixture = fixture;
}