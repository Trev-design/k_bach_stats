using Microsoft.EntityFrameworkCore;
using StackExchange.Redis;
using UserManager.Data;
using UserManager.Queries;
using UserManager.Rabbit;
using UserManager.Redis.Data;
using UserManager.Services;

var builder = WebApplication.CreateBuilder(args);

var redisConnString = builder.Configuration.GetConnectionString("RedisConnection") ?? throw new ArgumentException("connection does not exist");
var mySQLConnString = builder.Configuration.GetConnectionString("MySqlConnection") ?? throw new ArgumentException("connection does not exist");

builder.Services.AddDbContext<UserStoreContext>(options => options.UseMySql(
    mySQLConnString,
    new MySqlServerVersion(new Version(9, 0, 1))
));

builder.Services.AddSingleton<RabbitConn>();
builder.Services.AddSingleton<IConnectionMultiplexer>(opt => ConnectionMultiplexer.Connect(redisConnString));
builder.Services.AddScoped<ISessionRepo, RedisSessionRepo>();
builder.Services.AddHostedService<StartSessionService>();
builder.Services.AddHostedService<StopSessionService>();
builder.Services.AddHostedService<AddUserService>();
builder.Services.AddHostedService<DeleteUserService>();
builder.Services.AddGraphQLServer()
    .AddQueryType<UserQuery>();

var app = builder.Build();
app.MapGraphQL();
//app.UseHttpsRedirection();

Console.WriteLine("yeeeah");
app.Run();
Console.WriteLine("Ooooh");

