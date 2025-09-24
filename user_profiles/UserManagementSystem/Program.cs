using Microsoft.EntityFrameworkCore;
using UserManagementSystem.Services.Database;
using UserManagementSystem.Utils;

var builder = WebApplication.CreateBuilder(args);

builder.WebHost.ConfigureKestrel(options =>
{
    options.ListenAnyIP(5670, listenOptions =>
    {
        listenOptions.Protocols = Microsoft.AspNetCore.Server.Kestrel.Core.HttpProtocols.Http2;
    });
});

string connStr = DBConnstring.GetLocalEnvConnectionString();

builder.Services.AddDbContext<AppDBContext>(options => 
    options.UseMySql(connStr, 
                     ServerVersion.AutoDetect(connStr))
);

builder.Services.AddControllers();

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

app.MapControllers();

app.Run();
