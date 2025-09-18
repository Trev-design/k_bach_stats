
using Microsoft.EntityFrameworkCore;
using UserManagementSystem.Models;

namespace UserManagementSystem.Services.Database;

public class UserDBImpl
{
    public static async Task<List<User>> GetAllAsync(AppDBContext context)
    {
        var users = await context.Users.Include(user => user.UserProfile).
        ThenInclude(profile => profile.UserContact).
        ToListAsync();
        
        return users;
    }

    public static async Task CreateUser(AppDBContext context, string name, string email, string entity)
    {
        var userToInsert = new User { Entity = entity };
        await context.Users.AddAsync(userToInsert);
        var profileToInsert = new Profile { UserId = userToInsert.Id, ImagePath = "some_image.png", Description = "some description" };
        await context.Profiles.AddAsync(profileToInsert);
        var contactToInsert = new Contact { ProfileId = profileToInsert.Id, Email = email, Name = name };
        await context.Contacts.AddAsync(contactToInsert);
        await context.SaveChangesAsync();
    }
    public static async Task<User?> GetWholeUser(AppDBContext context, string entity)
    {
        var user = await context.Users.Include(u => u.UserProfile).
        ThenInclude(p => p.UserContact).
        Include(u => u.Contacts).
        Include(u => u.Workspaces).
        ThenInclude(w => w.ChatRooms).
        Include(u => u.Workspaces).
        ThenInclude(w => w.Contacts).
        FirstOrDefaultAsync(u => u.Entity == entity);

        return user;
    }
    public static async Task<User?> GetUserById(AppDBContext context, Guid id)
    {
        var user = await context.Users.Include(u => u.UserProfile).
        ThenInclude(p => p.UserContact).
        Include(u => u.Contacts).
        Include(u => u.Workspaces).
        FirstOrDefaultAsync(u => u.Id == id);

        return user;
    }

    public static async Task<Workspace?> AddNewWorkspace(AppDBContext context, Guid id, string workspaceName)
    {
        var user = await context.Users.FindAsync(id);
        if (user == null) return null;

        var workspace = new Workspace
        {
            UserId = id,
            Name = workspaceName
        };

        await context.Workspaces.AddAsync(workspace);
        await context.SaveChangesAsync();

        return workspace;
    }

    public static async Task DeleteWorkspace(AppDBContext context, Guid id, Guid workspaceId)
    {
        var workspace = await context.Workspaces.FirstOrDefaultAsync(w => w.UserId == id && w.Id == workspaceId) ?? throw new Exception("");
        context.Workspaces.Remove(workspace);
        await context.SaveChangesAsync();
    }

    public static async Task<ChatRoom?> NewChatRoom(AppDBContext context, Guid id, Guid workspaceId, string reference, string topic)
    {
        var workspace = await context.Workspaces.FirstOrDefaultAsync(w => w.Id == workspaceId && w.UserId == id);
        if (workspace == null) return null;

        var chatRoom = new ChatRoom
        {
            Reference = reference,
            WorkspaceId = workspaceId,
            Topic = topic
        };

        await context.ChatRooms.AddAsync(chatRoom);
        await context.SaveChangesAsync();

        return chatRoom;
    }

    public static async Task DeleteChat(AppDBContext context, Guid Id, Guid workspaceId, Guid chatId)
    {
        var workspace = await context.Workspaces.FirstOrDefaultAsync(w => w.Id == workspaceId && w.UserId == Id) ?? throw new Exception("");
        var chatRoom = await context.ChatRooms.FirstOrDefaultAsync(c => c.Id == chatId && c.WorkspaceId == workspaceId) ?? throw new Exception("");

        context.ChatRooms.Remove(chatRoom);
        await context.SaveChangesAsync();
    }

    public static async Task DeleteUser(AppDBContext context, Guid id)
    {
        var user = await context.Users.FirstOrDefaultAsync(u => u.Id == id) ?? throw new Exception("");
        context.Users.Remove(user);
        await context.SaveChangesAsync();
    }
}
