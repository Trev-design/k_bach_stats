using Microsoft.EntityFrameworkCore;
using UserManager.Data;
using UserManager.Models;

namespace UserManager.Queries;

public class UserQuery
{
    public async Task<Account> GetAccount([Service] UserStoreContext dbContext, [Service] HttpContextAccessor context)
    {
        Console.WriteLine("oooooohjaa");

        var entity = context.HttpContext?.Items["entity"] as string ?? throw new GraphQLException("invalid session");

        var account = await dbContext.Accounts
            .Include(account => account.AccountUser)
            .ThenInclude(user => user != null ? user.Profile : null)
            .ThenInclude(profile => profile != null ? profile.Contact : null)
            .Include(account => account.WorkSpaces)
            .FirstOrDefaultAsync(account => account.Entity == entity) ?? throw new GraphQLException("invalid entity");

        return account;
    }
}