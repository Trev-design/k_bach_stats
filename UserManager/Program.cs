using StackExchange.Redis;
using UserManager.Queries;
using UserManager.Redis.Data;
using UserManager.Serices;

var builder = WebApplication.CreateBuilder(args);

var connString = builder.Configuration.GetConnectionString("RedisConnection") ?? throw new ArgumentException("connection does not exist");

builder.Services.AddSingleton<IConnectionMultiplexer>(opt => ConnectionMultiplexer.Connect(connString));
builder.Services.AddSingleton<ISessionRepo, RedisSessionRepo>();
builder.Services.AddHostedService<RabbitConsumerService>();
builder.Services.AddGraphQLServer()
    .AddQueryType<UserQuery>();

var app = builder.Build();
app.MapGraphQL();
//app.UseHttpsRedirection();


app.Run();

