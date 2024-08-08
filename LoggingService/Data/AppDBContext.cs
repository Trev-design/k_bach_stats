using LoggingService.Models;
using Microsoft.EntityFrameworkCore;

namespace LoggingService.Data;

public class AppDBContext : DbContext
{
    public AppDBContext(DbContextOptions<AppDBContext> options) :base(options)
    {

    }

    public DbSet<Log> Logs => Set<Log>();
}