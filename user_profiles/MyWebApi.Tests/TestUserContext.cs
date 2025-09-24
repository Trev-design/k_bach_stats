using Microsoft.EntityFrameworkCore;
using Testcontainers.MySql;
using UserManagementSystem.Models;
using UserManagementSystem.Services.Database;

namespace MyWebApi.Tests;

[Collection("Database collection")]
public class TestUserContext(DatabaseFixture fixture)
{
    private readonly DatabaseFixture _fixture = fixture;

    [Fact]
    public async Task TestGetUserEntity()
    {
        var newUser = await GetNewUser("ajkamkalakjanhba@bjabhabah.jaj", "jajabababahba", "13");

        var user = await UserDBImpl.GetWholeUser(_fixture.Context, "13");
        Assert.NotNull(user);
        Assert.Equal("13", user.Entity);
    }

    [Fact]
    public async Task TestFetUserEntityFailed()
    {
        var newUser = await GetNewUser("acg@khaubaajb.ja", "lopoabjbahabhav", "12");

        var user = await UserDBImpl.GetWholeUser(_fixture.Context, "false_entity");
        Assert.Null(user);
    }

    [Fact]
    public async Task TestGetUser()
    {
        var newUser = await GetNewUser("havahjkahb@ahgag.hagh", "kanhjabhbakna", "11");

        var user = await UserDBImpl.GetUserById(_fixture.Context, newUser.Id);
        Assert.NotNull(user);
        Assert.Equal("11", user.Entity);
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
        var newUser = await GetNewUser("fadgfak@gahgakaj.kl", "trfagzuagtqkanu", "10");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, newUser.Id, "test_workspace");
        Assert.NotNull(workspace);

        var user = await UserDBImpl.GetUserById(_fixture.Context, workspace.UserId);
        Assert.NotNull(user);
        Assert.Equal("10", user.Entity);
    }

    [Fact]
    public async Task TestAddWorkspaceFailed()
    {
        var user = await GetNewUser("jjabgavha@lmakamka.ha", "hztaagazjaabzah", "9");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, Guid.NewGuid(), "test_workspace");
        Assert.Null(workspace);
    }

    [Fact]
    public async Task TestDeleteWorkspace()
    {
        var user = await GetNewUser("ghajkakj@gahgaj.hgaga", "jahajkak", "8");

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
        var user = await GetNewUser("fahja@hgah.ioio", "ghahjkalllkj", "7");

        await Assert.ThrowsAsync<Exception>(async () => await UserDBImpl.DeleteWorkspace(_fixture.Context, Guid.NewGuid(), Guid.NewGuid()));
    }

    [Fact]
    public async Task TestAddChat()
    {
        var user = await GetNewUser("jahaajah@hkak.klo", "fghajhak", "6");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, user.Id, "test_workspace");
        Assert.NotNull(workspace);

        var chat = await UserDBImpl.NewChatRoom(_fixture.Context, user.Id, workspace.Id, "Hello", "Mello");
        Assert.NotNull(chat);
    }

    [Fact]
    public async Task TestAddChatFailed()
    {
        var user = await GetNewUser("kazgauja@gah.ja", "jhabjbajanjan", "5");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, user.Id, "test_workspace");
        Assert.NotNull(workspace);

        var chat = await UserDBImpl.NewChatRoom(_fixture.Context, Guid.NewGuid(), workspace.Id, "Hello", "Mello");
        Assert.Null(chat);
    }

    [Fact]
    public async Task TestDeleteChat()
    {
        var user = await GetNewUser("hjkkjzfcvb@ah.kl", "kijzvcfvbgb", "4");

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
        var user = await GetNewUser("hikmk@hjkli.kj", "lmjhjhjkl", "3");

        var workspace = await UserDBImpl.AddNewWorkspace(_fixture.Context, user.Id, "test_workspace");
        Assert.NotNull(workspace);

        var chat = await UserDBImpl.NewChatRoom(_fixture.Context, user.Id, workspace.Id, "Hello", "Mello");
        Assert.NotNull(chat);

        await Assert.ThrowsAsync<Exception>(async () => await UserDBImpl.DeleteChat(_fixture.Context, user.Id, Guid.NewGuid(), chat.Id));
    }

    [Fact]
    public async Task TestDeleteUser()
    {
        var user = await GetNewUser("pkmn@jkl.lo", "kllmnb", "2");

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
        var user = await GetNewUser("ajla@ha.ql", "hgjhkl", "1");

        await Assert.ThrowsAsync<Exception>(async () => await UserDBImpl.DeleteUser(_fixture.Context, Guid.NewGuid()));
    }

    private async Task<User> GetNewUser(string email, string name, string entity)
    {
        var userToInsert = new User { Entity = entity };
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
