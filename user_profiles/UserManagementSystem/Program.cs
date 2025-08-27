using Microsoft.EntityFrameworkCore;
using UserManagementSystem.Services;

var builder = WebApplication.CreateBuilder(args);

var connStr = "Server=localhost;Database=kbach_users;User Id=kbach;Password=mysecretsqlpassword";

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
