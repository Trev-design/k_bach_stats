using System.Net.Http.Json;
using MyWebApi.Tests.Utils;
using UserManagementSystem.Models;

namespace MyWebApi.Tests;

[Collection("Database collection")]
public class TestUserController(EndpointsFixture fixture) : IClassFixture<EndpointsFixture>
{
    private readonly EndpointsFixture _fixture = fixture;

    [Fact]
    public async Task TestGetInitial()
    {
        foreach (var entity in _fixture.Entities)
        {
            var response = await _fixture.Client.GetAsync($"api/users/{entity}/initial");
            try
            {
                response.EnsureSuccessStatusCode();
                var created = await response.Content.ReadFromJsonAsync<User>() ?? throw new Exception("something wen wrong bei fetchin user");
            }
            catch (Exception e)
            {
                Assert.Fail($"Should succeed but got exception {e.Message}");
            }
        }
    }

    [Fact]
    public async Task TestGetInitialFailed()
    {
        await Assert.ThrowsAsync<HttpRequestException>(async () =>
        {
            var invalidID = Guid.NewGuid().ToString();
            var response = await _fixture.Client.GetAsync($"api/users/{invalidID}/initial");
            response.EnsureSuccessStatusCode();
        });
    }

    [Fact]
    public async Task TestGet()
    {
        foreach (var id in _fixture.UserIDs)
        {
            var response = await _fixture.Client.GetAsync($"api/users/{id}");
            try
            {
                response.EnsureSuccessStatusCode();
                var created = await response.Content.ReadFromJsonAsync<User>() ?? throw new Exception("something wen wrong bei fetchin user");
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
            var invalidID = Guid.NewGuid().ToString();
            var response = await _fixture.Client.GetAsync($"api/users/{invalidID}");
            response.EnsureSuccessStatusCode();
        });
    }

    [Fact]
    public async Task TestNewWorkspace()
    {
        foreach (var id in _fixture.UserIDs)
        {
            var userID = id.ToString();
            var name = RandomString.GenerateRandomString(20);
            var response = await _fixture.Client.PostAsJsonAsync($"api/users/{userID}/new_workspace", name);

            try
            {
                response.EnsureSuccessStatusCode();
                var created = await response.Content.ReadFromJsonAsync<Workspace>() ?? throw new Exception("could not fetch workspace");
                _fixture.Workspaces.Add(id, created.Id);
            }
            catch (Exception e)
            {
                Assert.Fail($"Should succeed but got exception {e.Message}");
            }
        }
    }

    [Fact]
    public async Task TestNewWorkspaceFailed()
    {
        await Assert.ThrowsAsync<HttpRequestException>(async () =>
        {
            var invalidUserID = Guid.NewGuid().ToString();
            var name = RandomString.GenerateRandomString(20);
            var response = await _fixture.Client.PostAsJsonAsync($"api/users/{invalidUserID}/new_workspace", name);
            response.EnsureSuccessStatusCode();
        });
    }

    [Fact]
    public async Task TestDeleteWorkspace()
    {
        foreach (var id in _fixture.UserIDs)
        {
            bool ok = _fixture.Workspaces.TryGetValue(id, out Guid workspaceID);
            Assert.True(ok);

            var response = await _fixture.Client.DeleteAsync($"api/users/{id.ToString()}/workspace/{workspaceID.ToString()}");
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
    public async Task TestDeleteWorkspaceFailed()
    {
        await Assert.ThrowsAsync<HttpRequestException>(async () =>
        {
            var invalidUserID = Guid.NewGuid().ToString();
            var invalidWorkspaceID = Guid.NewGuid().ToString();
            var response = await _fixture.Client.DeleteAsync($"api/users/{invalidUserID}/workspace/{invalidWorkspaceID}");
            response.EnsureSuccessStatusCode();
        });
    }

    [Fact]
    public async Task TestNewChat()
    {
        foreach (var id in _fixture.UserIDs)
        {
            bool ok = _fixture.Workspaces.TryGetValue(id, out Guid workspaceID);
            Assert.True(ok);
            var topic = RandomString.GenerateRandomString(20);
            var response = await _fixture.Client.PostAsJsonAsync($"api/users/{id.ToString()}/workspace/{workspaceID.ToString()}/new_chat", topic);

            try
            {
                response.EnsureSuccessStatusCode();
                var created = await response.Content.ReadFromJsonAsync<ChatRoom>() ?? throw new Exception("could not fetch chatroom from database");
                _fixture.ChatRooms.Add(id, created.Id);
            }
            catch (Exception e)
            {
                Assert.Fail($"Should succeed but got exception {e.Message}");
            }
        }
    }

    [Fact]
    public async Task TestNewChatFailed()
    {
        await Assert.ThrowsAsync<HttpRequestException>(async () =>
        {
            var invalidUserID = Guid.NewGuid().ToString();
            var invalidWorkspaceID = Guid.NewGuid().ToString();
            var topic = RandomString.GenerateRandomString(20);
            var response = await _fixture.Client.PostAsJsonAsync($"api/users/{invalidUserID}/workspace/{invalidWorkspaceID}/new_chat", topic);
            response.EnsureSuccessStatusCode();
        });
    }

    [Fact]
    public async Task TestDeleteChat()
    {
        foreach (var id in _fixture.UserIDs)
        {
            bool ok = _fixture.Workspaces.TryGetValue(id, out Guid workspaceID);
            Assert.True(ok);
            ok = _fixture.ChatRooms.TryGetValue(workspaceID, out Guid chatroomID);
            Assert.True(ok);

            var response = await _fixture.Client.DeleteAsync($"api/users/{id.ToString()}/workspace/{workspaceID.ToString()}/chat/{chatroomID.ToString()}");
            try
            {
                response.EnsureSuccessStatusCode();
                var created = await response.Content.ReadFromJsonAsync<User>() ?? throw new Exception("something wen wrong bei fetchin user");
            }
            catch (Exception e)
            {
                Assert.Fail($"Should succeed but got exception {e.Message}");
            }
        }
    }

    [Fact]
    public async Task TestDeleteChatroomFailed()
    {
        await Assert.ThrowsAsync<HttpRequestException>(async () =>
        {
            var invalidUserID = Guid.NewGuid().ToString();
            var invalidProfileID = Guid.NewGuid().ToString();
            var invalidChatroomID = Guid.NewGuid().ToString();

            var response = await _fixture.Client.DeleteAsync($"api/users/{invalidUserID}/workspace/{invalidProfileID}/chat/{invalidChatroomID}");
            response.EnsureSuccessStatusCode();
        });
    }
}

