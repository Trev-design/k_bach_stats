using System.Security.AccessControl;
using Microsoft.EntityFrameworkCore;
using UserManager.Models;

namespace UserManager.Data;

public class UserStoreContext(DbContextOptions<UserStoreContext> options) : DbContext(options) 
{
    public DbSet<Account> Accounts => Set<Account>();
    public DbSet<User> Users => Set<User>();
    public DbSet<Profile> Profiles => Set<Profile>();
    public DbSet<Contact> Contacts => Set<Contact>();
    public DbSet<WorkSpace> WorkSpaces => Set<WorkSpace>();
    public DbSet<Experience> Experiences => Set<Experience>();
    public DbSet<SelfAssessment> SelfAssessments => Set<SelfAssessment>();
}