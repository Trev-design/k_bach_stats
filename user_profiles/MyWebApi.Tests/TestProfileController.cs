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
        foreach (var id in _fixture.ProfileIDs)
        {
            var response = await _fixture.Client.GetAsync($"api/profile/{id}");
            try
            {
                response.EnsureSuccessStatusCode();
                var created = await response.Content.ReadFromJsonAsync<Profile>() ?? throw new Exception("something wen wrong bei fetchin profile");
            }
            catch (Exception e)
            {
                Assert.Fail($"Should succeed but got exception {e.Message}");
            }
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
        foreach (var id in _fixture.ProfileIDs)
        {
            var response = await _fixture.Client.PutAsJsonAsync($"api/profile/{id}/new_image", RandomString.GenerateRandomString(100));
            try
            {
                response.EnsureSuccessStatusCode();
            }
            catch (Exception e)
            {
                Assert.Fail($"Should succeed but got exception {e.Message}");
            }
        }
    }

    [Fact]
    public async Task TestChangeImageFailed()
    {
        await Assert.ThrowsAsync<HttpRequestException>(async () =>
        {
            var id = Guid.NewGuid();
            var response = await _fixture.Client.PutAsJsonAsync($"api/profile/{id}", RandomString.GenerateRandomString(150));
            response.EnsureSuccessStatusCode();
        });
    }

    [Fact]
    public async Task TestChangeDescription()
    {
        foreach (var id in _fixture.ProfileIDs)
        {
            var response = await _fixture.Client.PutAsJsonAsync($"api/profile/{id}/new_image", RandomString.GenerateRandomString(100));
            try
            {
                response.EnsureSuccessStatusCode();
            }
            catch (Exception e)
            {
                Assert.Fail($"Should succeed but got exception {e.Message}");
            }
        }
    }

    [Fact]
    public async Task TestChangeDescriptionFailed()
    {
        await Assert.ThrowsAsync<HttpRequestException>(async () =>
        {
            var id = Guid.NewGuid();
            var response = await _fixture.Client.GetAsync($"api/profile/{id}");
            response.EnsureSuccessStatusCode();
        });
    }
}