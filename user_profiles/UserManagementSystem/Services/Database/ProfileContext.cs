using Microsoft.EntityFrameworkCore;
using UserManagementSystem.Models;

namespace UserManagementSystem.Services.Database;

/// <summary>
/// implements access for the profiles
/// </summary>
public class ProfileDBImpl
{
    /// <summary>
    /// gets a profile that includes some contact information
    /// </summary>
    /// <param name="context"></param>
    /// <param name="id"></param>
    /// <returns>A Profile on the gien id</returns>
    public static async Task<Profile?> GetProfile(AppDBContext context, Guid id)
    {
        var profile = await context.Profiles.Include(p => p.UserContact).FirstOrDefaultAsync(p => p.Id == id);
        if (profile == null) return null;
        return profile;
    }

    /// <summary>
    /// to make it possible to change the profile image
    /// </summary>
    /// <param name="context"></param>
    /// <param name="id"></param>
    /// <param name="imagePath"></param>
    /// <returns></returns>
    /// <exception cref="Exception"></exception>
    public static async Task ChangeImage(AppDBContext context, Guid id, string imagePath)
    {
        var profile = await context.Profiles.FirstOrDefaultAsync(p => p.Id == id) ?? throw new Exception("");
        profile.ImagePath = imagePath;
        await context.SaveChangesAsync();
    }

    /// <summary>
    /// to make it possible to change the profile description
    /// </summary>
    /// <param name="context"></param>
    /// <param name="id"></param>
    /// <param name="description"></param>
    /// <returns></returns>
    /// <exception cref="Exception"></exception>
    public static async Task ChangeDescription(AppDBContext context, Guid id, string description)
    {
        var profile = await context.Profiles.FirstOrDefaultAsync(p => p.Id == id) ?? throw new Exception(""); ;
        profile.Description = description;
        await context.SaveChangesAsync();
    }

    /// <summary>
    /// to make itposible to change the name
    /// </summary>
    /// <param name="context"></param>
    /// <param name="id"></param>
    /// <param name="contactId"></param>
    /// <param name="contactName"></param>
    /// <returns></returns>
    /// <exception cref="Exception"></exception>
    public static async Task ChangeContactName(AppDBContext context, Guid id, Guid contactId, string contactName)
    {
        var contact = await context.Contacts.FirstOrDefaultAsync(c => c.Id == contactId && c.ProfileId == id) ?? throw new Exception("");
        contact.Name = contactName;
        await context.SaveChangesAsync();
    }

    /// <summary>
    /// make it possible to change the email
    /// </summary>
    /// <param name="context"></param>
    /// <param name="id"></param>
    /// <param name="contactId"></param>
    /// <param name="email"></param>
    /// <returns></returns>
    /// <exception cref="Exception"></exception>
    public static async Task ChangeContactEmail(AppDBContext context, Guid id, Guid contactId, string email)
    {
        var contact = await context.Contacts.FirstOrDefaultAsync(c => c.Id == contactId && c.ProfileId == id) ?? throw new Exception("");
        contact.Email = email;
        await context.SaveChangesAsync();
    }
}