using Microsoft.EntityFrameworkCore;
using UserManagementSystem.Models;

namespace UserManagementSystem.Services;

public class ProfileDBImpl
{
    public static async Task<Profile?> GetProfile(AppDBContext context, Guid id)
    {
        var profile = await context.Profiles.Include(p => p.UserContact).FirstOrDefaultAsync(p => p.Id == id);
        if (profile == null) return null;
        return profile;
    }

    public static async Task ChangeImage(AppDBContext context, Guid id, string imagePath)
    {
        var profile = await context.Profiles.FirstOrDefaultAsync(p => p.Id == id) ?? throw new Exception("");
        profile.ImagePath = imagePath;
        await context.SaveChangesAsync();
    }

    public static async Task ChangeDescription(AppDBContext context, Guid id, string description)
    {
        var profile = await context.Profiles.FirstOrDefaultAsync(p => p.Id == id) ?? throw new Exception(""); ;
        profile.Description = description;
        await context.SaveChangesAsync();
    }

    public static async Task ChangeContactName(AppDBContext context, Guid id, Guid contactId, string contactName)
    {
        var contact = await context.Contacts.FirstOrDefaultAsync(c => c.Id == contactId && c.ProfileId == id) ?? throw new Exception("");
        contact.Name = contactName;
        await context.SaveChangesAsync();
    }

    public static async Task ChangeContactEmail(AppDBContext context, Guid id, Guid contactId, string email)
    {
        var contact = await context.Contacts.FirstOrDefaultAsync(c => c.Id == contactId && c.ProfileId == id) ?? throw new Exception("");
        contact.Email = email;
        await context.SaveChangesAsync();
    }
}