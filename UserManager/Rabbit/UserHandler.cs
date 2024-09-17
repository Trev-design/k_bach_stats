using System.Text;
using System.Text.Json;
using Microsoft.EntityFrameworkCore;
using UserManager.Data;
using UserManager.Models;

namespace UserManager.Rabbit;

public class UserHandler(IServiceScopeFactory scopeFactory)
{
    private readonly IServiceScopeFactory _scopeFactory = scopeFactory;

    public async Task MakeUser(ReadOnlyMemory<byte> body)
    {
        using var scope = _scopeFactory.CreateAsyncScope();
        var context = scope.ServiceProvider.GetRequiredService<UserStoreContext>();

        string jsonData = Encoding.UTF8.GetString(body.Span);
        AccountData data = JsonSerializer.Deserialize<AccountData>(jsonData) ?? throw new ArgumentException("invalid credentials");

        var newAccount = new Account { Entity = data.Entity };

        var newUser = new User { AccountId = newAccount.Id };

        var newProfile = new Profile { UserId = newUser.Id };

        var newContact = new Contact { ProfileId = newProfile.Id, Name = data.UserName, Email = data.Email };

        newAccount.AccountUser = newUser;
        newUser.Profile = newProfile;
        newProfile.Contact = newContact;

        await context.Accounts.AddAsync(newAccount);
        await context.SaveChangesAsync();
    }

    public async Task DeleteUser(ReadOnlyMemory<byte> body)
    {
        using var scope = _scopeFactory.CreateAsyncScope();
        var context = scope.ServiceProvider.GetRequiredService<UserStoreContext>();

        string jsonData = Encoding.UTF8.GetString(body.Span);
        DeleteAccountData accountToDelete = JsonSerializer.Deserialize<DeleteAccountData>(jsonData) ?? throw new ArgumentException("invalid data");

        Account? account;

        try {
            account = 
                await context
                  .Accounts
                  .Where(acc => acc.Entity == accountToDelete.Id.ToString())
                  .SingleOrDefaultAsync();
        } 
        catch (Exception ex) {
            throw new ArgumentException(ex.Message);
        }

        if (account == null) 
        {
            throw new ArgumentException("invalid account data");
        }

        context.Accounts.Remove(account);
        await context.SaveChangesAsync();
    }
}