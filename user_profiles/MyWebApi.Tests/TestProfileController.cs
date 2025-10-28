using System.Net.Http;
using System.Net.Http.Json;
using MyWebApi.Tests.Utils;
using UserManagementSystem.Models;

namespace MyWebApi.Tests;

[Collection("Database collection")]
public class TestProfileController(EndpointsFixture fixture) : IClassFixture<EndpointsFixture>
{
    private readonly EndpointsFixture _fixture = fixture;

    [Fact]
    public async Task TestGet()
    {
        foreach (var id in _fixture.UserIDs)
        {
            var getUserResponse = await _fixture.Client.GetAsync($"/api/users/{id}");

            try
            {
                getUserResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var user = await getUserResponse.Content.ReadFromJsonAsync<User>();
            Assert.NotNull(user);

            var response = await _fixture.Client.GetAsync($"api/profile/{user.UserProfile.Id}");

            try
            {
                getUserResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var profile = await response.Content.ReadFromJsonAsync<Profile>();
            Assert.NotNull(profile);
        }
    }

    [Fact]
    public async Task TestGetFailed()
    {
        await Assert.ThrowsAsync<HttpRequestException>(async () =>
        {
            var id = Guid.NewGuid();
            var response = await _fixture.Client.GetAsync($"api/profile/{id}");
            response.EnsureSuccessStatusCode();
        });
    }

    [Fact]
    public async Task TestChangeImage()
    {
        foreach (var id in _fixture.UserIDs)
        {
            var getUserResponse = await _fixture.Client.GetAsync($"api/users/{id}");

            try
            {
                getUserResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var user = await getUserResponse.Content.ReadFromJsonAsync<User>();
            Assert.NotNull(user);

            var profile = user.UserProfile;

            var newImagePath = RandomString.GenerateRandomString(32);

            var response = await _fixture.Client.PutAsJsonAsync($"api/profile/{profile.Id}/new_image", newImagePath);

            try
            {
                response.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var getChangedProfileResponse = await _fixture.Client.GetAsync($"api/profile/{profile.Id}");

            try
            {
                getChangedProfileResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var newProfile = await getChangedProfileResponse.Content.ReadFromJsonAsync<Profile>();
            Assert.NotNull(newProfile);
            Assert.Equal(newImagePath, newProfile.ImagePath);
        }
    }

    [Fact]
    public async Task TestChangeImageFailed()
    {
        await Assert.ThrowsAsync<HttpRequestException>(async () =>
        {
            var id = Guid.NewGuid();
            var response = await _fixture.Client.PutAsJsonAsync($"api/profile/{id}/new_image", RandomString.GenerateRandomString(150));
            response.EnsureSuccessStatusCode();
        });
    }

    [Fact]
    public async Task TestChangeDescription()
    {
        foreach (var id in _fixture.UserIDs)
        {
            var getUserResponse = await _fixture.Client.GetAsync($"api/users/{id}");

            try
            {
                getUserResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var user = await getUserResponse.Content.ReadFromJsonAsync<User>();
            Assert.NotNull(user);

            var profile = user.UserProfile;

            var newDescription = RandomString.GenerateRandomString(32);

            var response = await _fixture.Client.PutAsJsonAsync($"api/profile/{profile.Id}/new_description", newDescription);

            try
            {
                response.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var getChangedProfileResponse = await _fixture.Client.GetAsync($"api/profile/{profile.Id}");

            try
            {
                getChangedProfileResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var newProfile = await getChangedProfileResponse.Content.ReadFromJsonAsync<Profile>();
            Assert.NotNull(newProfile);
            Assert.Equal(newDescription, newProfile.Description);
        }
    }

    [Fact]
    public async Task TestChangeDescriptionFailed()
    {
        await Assert.ThrowsAsync<HttpRequestException>(async () =>
        {
            var id = Guid.NewGuid();
            var response = await _fixture.Client.GetAsync($"api/profile/{id}/new_description");
            response.EnsureSuccessStatusCode();
        });
    }
}