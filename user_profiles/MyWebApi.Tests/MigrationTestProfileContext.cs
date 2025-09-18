
using Microsoft.EntityFrameworkCore;
using Testcontainers.MySql;
using UserManagementSystem.Models;
using UserManagementSystem.Services.Database;

namespace MyWebApi.Tests;

[Collection("Database collection")]
public class MigrationTestProfileContext(DatabaseFixture fixture)
{
    private readonly DatabaseFixture _fixture = fixture;

    [Fact]
    public async Task TestGetProfile()
    {
        var newProfile = await MakeProfile("pp@aa.ss", "23");
        Assert.NotNull(newProfile);

        var profile = await ProfileDBImpl.GetProfile(_fixture.Context, newProfile.Id);
        Assert.NotNull(profile);
    }

    [Fact]
    public async Task TestGetProfileFailed()
    {
        var profile = await ProfileDBImpl.GetProfile(_fixture.Context, Guid.NewGuid());
        Assert.Null(profile);
    }

    [Fact]
    public async Task TestChangeDescription()
    {
        var newProfile = await MakeProfile("uu@ii.oo", "22");
        Assert.NotNull(newProfile);

        try
        {
            await ProfileDBImpl.ChangeDescription(_fixture.Context, newProfile.Id, "new_description");
        }
        catch
        {
            Assert.Fail("should be successfull but got an exception");
        }

        var profile = await ProfileDBImpl.GetProfile(_fixture.Context, newProfile.Id);
        Assert.NotNull(profile);
        Assert.Equal("new_description", profile.Description);
    }

    [Fact]
    public async Task TestChangeDescriptionFailed()
    {
        await Assert.ThrowsAsync<Exception>(async () => await ProfileDBImpl.ChangeDescription(_fixture.Context, Guid.NewGuid(), "aaaaaaaaah"));
    }

    [Fact]
    public async Task TestChangeContactName()
    {
        var newProfile = await MakeProfile("rr@tt.zz", "21");
        Assert.NotNull(newProfile);

        try
        {
            await ProfileDBImpl.ChangeContactName(_fixture.Context, newProfile.Id, newProfile.UserContact.Id, "NewName");
        }
        catch
        {
            Assert.Fail("should be successfull but got an exception");
        }

        var profile = await ProfileDBImpl.GetProfile(_fixture.Context, newProfile.Id);
        Assert.NotNull(profile);
        Assert.Equal("NewName", profile.UserContact.Name);
    }

    [Fact]
    public async Task TestChangeContactNameFailed()
    {
        await Assert.ThrowsAsync<Exception>(async () => await ProfileDBImpl.ChangeContactName(_fixture.Context, Guid.NewGuid(), Guid.NewGuid(), ""));
    }

    [Fact]
    public async Task TestChangeContactEmail()
    {
        var newProfile = await MakeProfile("qq@ww.ee", "20");
        Assert.NotNull(newProfile);

        try
        {
            await ProfileDBImpl.ChangeContactEmail(_fixture.Context, newProfile.Id, newProfile.UserContact.Id, "new@email.com");
        }
        catch
        {
            Assert.Fail("should be successfull but got an exception");
        }

        var profile = await ProfileDBImpl.GetProfile(_fixture.Context, newProfile.Id);
        Assert.NotNull(profile);
        Assert.Equal("new@email.com", profile.UserContact.Email);
    }

    [Fact]
    public async Task TestChangeContactEmailFailed()
    {
        await Assert.ThrowsAsync<Exception>(async () =>
            await ProfileDBImpl.ChangeContactEmail(
                _fixture.Context,
                Guid.NewGuid(),
                Guid.NewGuid(),
                "aaaaah"
        ));
    }

    private async Task<Profile> MakeProfile(string email, string entity)
    {
        var user = new User { Entity = entity };
        await _fixture.Context.Users.AddAsync(user);
        var profile = new Profile { UserId = user.Id, ImagePath = "some_image.jpg", Description = "some description", ProfileUser = user };
        await _fixture.Context.Profiles.AddAsync(profile);
        var contact = new Contact { Name = "MyName", Email = email, ProfileId = profile.Id, UserProfile = profile };
        await _fixture.Context.Contacts.AddAsync(contact);
        await _fixture.Context.SaveChangesAsync();
        profile.UserContact = contact;
        return profile;
    }
}