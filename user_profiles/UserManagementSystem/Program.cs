using Microsoft.EntityFrameworkCore;
using Minio;
using Minio.DataModel.Args;
using UserManagementSystem.Services.Database;
using UserManagementSystem.Services.RabbitMQ;
using UserManagementSystem.Services.S3Service;
using UserManagementSystem.Utils;

var builder = WebApplication.CreateBuilder(args);

builder.WebHost.ConfigureKestrel(options =>
{
    options.ListenAnyIP(5670, listenOptions =>
    {
        listenOptions.Protocols = Microsoft.AspNetCore.Server.Kestrel.Core.HttpProtocols.Http2;
    });
});

builder.Services.AddSingleton<S3Settings>();

builder.Services.AddSingleton(async (provider) =>
{
    var settings = provider.GetRequiredService<S3Settings>();

    return new MinioClient()
        .WithEndpoint(settings.URL)
        .WithCredentials(settings.GetUser(), settings.GetSekretKey())
        .Build();
});

builder.Services.AddSingleton<S3Handler>();

builder.Services.AddSingleton<IMessageChannel, RabbitMessageChannel>();
builder.Services.AddHostedService<RabbitMQLoggingService>();

string connStr = DBConnstring.GetLocalEnvConnectionString();

builder.Services.AddDbContext<AppDBContext>(options => 
    options.UseMySql(connStr, 
                     ServerVersion.AutoDetect(connStr))
);

builder.Services.AddControllers();

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

using (var scope = app.Services.CreateScope())
{
    var minio = scope.ServiceProvider.GetRequiredService<IMinioClient>();
    var settings = scope.ServiceProvider.GetRequiredService<S3Settings>();

    bool found = await minio.BucketExistsAsync(new BucketExistsArgs().WithBucket(settings.BucketName));
    if (!found) await minio.MakeBucketAsync(new MakeBucketArgs().WithBucket(settings.BucketName));
}

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

app.MapControllers();

app.Run();
