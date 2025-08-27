using Microsoft.EntityFrameworkCore;
using Testcontainers.MySql;
using UserManagementSystem.Models;
using UserManagementSystem.Services;

namespace MyWebApi.Tests;

[Collection("Database collection")]
public class MigrationTestUserContext(DatabaseFixture fixture)
{
    private readonly DatabaseFixture _fixture = fixture;

    [Fact]
    public async Task TestGetUserEntity()
    {
        var newUser = await GetNewUser("ajkamkalakjanhba@bjabhabah.jaj", "jajabababahba");

        var user = await UserDBImpl.GetWholeUser(_fixture.Context, "entity");
        Assert.NotNull(user);
        Assert.Equal("entity", user.Entity);
    }

    [Fact]
    public async Task TestFetUserEntityFailed()
    {
        var newUser = await GetNewUser("acg@khaubaajb.ja", "lopoabjbahabhav");

        var user = await UserDBImpl.GetWholeUser(_fixture.Context, "false_entity");
        Assert.Null(user);
    }

    [Fact]
    public async Task TestGetUser()
    {
        var newUser = await GetNewUser("havahjkahb@ahgag.hagh", "kanhjabhbakna");

        var user = await UserDBImpl.GetUserById(_fixture.Context, newUser.Id);
        Assert.NotNull(user);
        Assert.Equal("entity", user.Entity);
    }

    [Fact]
    public async Task TestGetUserFailed()
    {
        var user = await UserDBImpl.GetUserById(_fixture.Context, Guid.NewGuid());
        Assert.Null(user);
    }

    [Fact]
    public async Task TestAddWorkspace()
    {
        var newUser = await GetNewUser("fadgfak@gahgakaj.kl", "trfagzuagtqkanu");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, newUser.Id, "test_workspace");
        Assert.NotNull(workspace);

        var user = await UserDBImpl.GetUserById(_fixture.Context, workspace.UserId);
        Assert.NotNull(user);
        Assert.Equal("entity", user.Entity);
    }

    [Fact]
    public async Task TestAddWorkspaceFailed()
    {
        var user = await GetNewUser("jjabgavha@lmakamka.ha", "hztaagazjaabzah");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, Guid.NewGuid(), "test_workspace");
        Assert.Null(workspace);
    }

    [Fact]
    public async Task TestDeleteWorkspace()
    {
        var user = await GetNewUser("ghajkakj@gahgaj.hgaga", "jahajkak");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, user.Id, "test_workspace");
        Assert.NotNull(workspace);

        try
        {
            await UserDBImpl.DeleteWorkspace(_fixture.Context, user.Id, workspace.Id);
        }
        catch
        {
            Assert.Fail("should be successfull but got an exception");
        }

        var nullWorkspace = _fixture.Context.Workspaces.FirstOrDefault(w => w.Id == workspace.Id);
        Assert.Null(nullWorkspace);
    }

    [Fact]
    public async Task TestDeleteWorkspaceFailed()
    {
        var user = await GetNewUser("fahja@hgah.ioio", "ghahjkalllkj");

        await Assert.ThrowsAsync<Exception>(async () => await UserDBImpl.DeleteWorkspace(_fixture.Context, Guid.NewGuid(), Guid.NewGuid()));
    }

    [Fact]
    public async Task TestAddChat()
    {
        var user = await GetNewUser("jahaajah@hkak.klo", "fghajhak");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, user.Id, "test_workspace");
        Assert.NotNull(workspace);

        var chat = await UserDBImpl.NewChatRoom(_fixture.Context, user.Id, workspace.Id, "Hello", "Mello");
        Assert.NotNull(chat);
    }

    [Fact]
    public async Task TestAddChatFailed()
    {
        var user = await GetNewUser("kazgauja@gah.ja", "jhabjbajanjan");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, user.Id, "test_workspace");
        Assert.NotNull(workspace);

        var chat = await UserDBImpl.NewChatRoom(_fixture.Context, Guid.NewGuid(), workspace.Id, "Hello", "Mello");
        Assert.Null(chat);
    }

    [Fact]
    public async Task TestDeleteChat()
    {
        var user = await GetNewUser("hjkkjzfcvb@ah.kl", "kijzvcfvbgb");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, user.Id, "test_workspace");
        Assert.NotNull(workspace);

        var chat = await UserDBImpl.NewChatRoom(_fixture.Context, user.Id, workspace.Id, "Hello", "Mello");
        Assert.NotNull(chat);

        try
        {
            await UserDBImpl.DeleteChat(_fixture.Context, user.Id, workspace.Id, chat.Id);
        }
        catch
        {
            Assert.Fail("should be successfull but got an exception");
        }

        var nullChat = await _fixture.Context.ChatRooms.FirstOrDefaultAsync(c => c.Id == chat.Id);
        Assert.Null(nullChat);
    }

    [Fact]
    public async Task TestDeleteChatFailed()
    {
        var user = await GetNewUser("hikmk@hjkli.kj", "lmjhjhjkl");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, user.Id, "test_workspace");
        Assert.NotNull(workspace);

        var chat = await UserDBImpl.NewChatRoom(_fixture.Context, user.Id, workspace.Id, "Hello", "Mello");
        Assert.NotNull(chat);

        await Assert.ThrowsAsync<Exception>(async () => await UserDBImpl.DeleteChat(_fixture.Context, user.Id, Guid.NewGuid(), chat.Id));
    }

    [Fact]
    public async Task TestDeleteUser()
    {
        var user = await GetNewUser("pkmn@jkl.lo", "kllmnb");

        try
        {
            await UserDBImpl.DeleteUser(_fixture.Context, user.Id);
        }
        catch
        {
            Assert.Fail("should be successfull but got an exception");
        }

        var nullUser = await UserDBImpl.GetUserById(_fixture.Context, user.Id);
        Assert.Null(nullUser);
    }

    [Fact]
    public async Task TestDeletUserFailed()
    {
        var user = await GetNewUser("ajla@ha.ql", "hgjhkl");

        await Assert.ThrowsAsync<Exception>(async () => await UserDBImpl.DeleteUser(_fixture.Context, Guid.NewGuid()));
    }

    private async Task<User> GetNewUser(string email, string name)
    {
        var userToInsert = new User { Entity = "entity" };
        await _fixture.Context.Users.AddAsync(userToInsert);
        var profileToInsert = new Profile { UserId = userToInsert.Id, ImagePath = "some_image.png", Description = "some description" };
        await _fixture.Context.Profiles.AddAsync(profileToInsert);
        var contactToInsert = new Contact { ProfileId = profileToInsert.Id, Email = email, Name = name };
        await _fixture.Context.Contacts.AddAsync(contactToInsert);
        await _fixture.Context.SaveChangesAsync();

        profileToInsert.UserContact = contactToInsert;
        userToInsert.UserProfile = profileToInsert;

        return userToInsert;
    }
}
