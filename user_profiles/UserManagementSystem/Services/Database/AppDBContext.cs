using Microsoft.EntityFrameworkCore;
using UserManagementSystem.Models;

namespace UserManagementSystem.Services.Database;

public class AppDBContext(DbContextOptions<AppDBContext> options) : DbContext(options)
{
    public DbSet<User> Users => Set<User>();
    public DbSet<Profile> Profiles => Set<Profile>();
    public DbSet<Contact> Contacts => Set<Contact>();
    public DbSet<Workspace> Workspaces => Set<Workspace>();
    public DbSet<ChatRoom> ChatRooms => Set<ChatRoom>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        // generate the ids on add
        modelBuilder.Entity<User>().Property(u => u.Id).ValueGeneratedOnAdd();
        modelBuilder.Entity<Profile>().Property(p => p.Id).ValueGeneratedOnAdd();
        modelBuilder.Entity<Contact>().Property(c => c.Id).ValueGeneratedOnAdd();
        modelBuilder.Entity<Workspace>().Property(ws => ws.Id).ValueGeneratedOnAdd();
        modelBuilder.Entity<ChatRoom>().Property(cr => cr.Id).ValueGeneratedOnAdd();

        // User -> Profile (1:1) mit Cascade Delete
        modelBuilder.Entity<User>().
            HasOne(u => u.UserProfile).
            WithOne(p => p.ProfileUser).
            HasForeignKey<Profile>(p => p.UserId).
            OnDelete(DeleteBehavior.Cascade);

        // Profile -> Contact (1:1) mit Cascade Delete
        modelBuilder.Entity<Profile>().
            HasOne(p => p.UserContact).
            WithOne(c => c.UserProfile).
            HasForeignKey<Contact>(c => c.ProfileId).
            OnDelete(DeleteBehavior.Cascade);

        // User -> Workspaces (1:n) mit Cascade Delete
        modelBuilder.Entity<User>().
            HasMany(u => u.Workspaces).
            WithOne(w => w.User).
            HasForeignKey(w => w.UserId).
            OnDelete(DeleteBehavior.Cascade);

        // User -> Contacts (Many-to-Many)
        modelBuilder.Entity<User>().
            HasMany(u => u.Contacts).
            WithMany(c => c.Users).
            UsingEntity(j => j.ToTable("UserContacts"));

        // Workspace -> Contacts (Many-to-Many)
        modelBuilder.Entity<Workspace>().
            HasMany(w => w.Contacts).
            WithMany(c => c.Workspaces).
            UsingEntity(j => j.ToTable("WorkspaceContacts"));

        // Workspace -> ChatRooms (1:n) mit Cascade Delete
        modelBuilder.Entity<Workspace>().
            HasMany(w => w.ChatRooms).
            WithOne(c => c.Workspace).
            HasForeignKey(c => c.WorkspaceId).
            OnDelete(DeleteBehavior.Cascade);
    }
}