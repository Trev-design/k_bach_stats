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
        foreach (var entity in _fixture.EntitiesWithChatRooms)
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
            var name = RandomString.GenerateRandomString(20);
            var response = await _fixture.Client.PostAsJsonAsync($"api/users/{id}/new_workspace", name);

            try
            {
                response.EnsureSuccessStatusCode();
            }
            catch (Exception e)
            {
                Assert.Fail($"Should succeed but got exception {e.Message}");
            }

            var created = await response.Content.ReadFromJsonAsync<Workspace>();
            Assert.NotNull(created);
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
        foreach (var id in _fixture.DeleteWorcspaceIDs)
        {
            var getResponse = await _fixture.Client.GetAsync($"api/users/{id}");

            try
            {
                getResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var user = await getResponse.Content.ReadFromJsonAsync<User>();
            Assert.NotNull(user);

            foreach (var workspace in user.Workspaces)
            {
                var response = await _fixture.Client.DeleteAsync($"api/users/{user.Id}/workspace/{workspace.Id}");
                try
                {
                    response.EnsureSuccessStatusCode();
                }
                catch (HttpRequestException e)
                {
                    Assert.Fail($"Should succeed but got exception {e.Message}");
                }
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
        foreach (var id in _fixture.UserIDsWithWorkspaces)
        {
            var getResponse = await _fixture.Client.GetAsync($"api/users/{id}");
            try
            {
                getResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var user = await getResponse.Content.ReadFromJsonAsync<User>();
            Assert.NotNull(user);

            foreach (var workspace in user.Workspaces)
            {
                var response = await _fixture.Client.PostAsJsonAsync($"api/users/{user.Id}/workspace/{workspace.Id}/new_chat", RandomString.GenerateRandomString(40));

                try
                {
                    response.EnsureSuccessStatusCode();
                }
                catch (HttpRequestException exception)
                {
                    Assert.Fail(exception.Message);
                }

                var created = await response.Content.ReadFromJsonAsync<ChatRoom>();
                Assert.NotNull(created);
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
        foreach (var id in _fixture.DeleteChatIDs)
        {
            var getResponse = await _fixture.Client.GetAsync($"api/users/{id}");

            try
            {
                getResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var user = await getResponse.Content.ReadFromJsonAsync<User>();
            Assert.NotNull(user);
            
            foreach (var workspace in user.Workspaces)
            {
                foreach (var chatroom in workspace.ChatRooms)
                {
                    var response = await _fixture.Client.DeleteAsync($"api/users/{user.Id}/workspace/{workspace.Id}/chat/{chatroom.Id}");

                    try
                    {
                        response.EnsureSuccessStatusCode();
                    } catch (HttpRequestException exception)
                    {
                        Assert.Fail(exception.Message);
                    }
                }
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

    [Fact]
    public async Task TestAddWorkspaceByEntity()
    {
        foreach (var entity in _fixture.Entities)
        {
            var getResponse = await _fixture.Client.GetAsync($"api/users/{entity}/initial");

            try
            {
                getResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var user = await getResponse.Content.ReadFromJsonAsync<User>();
            Assert.NotNull(user);

            var name = RandomString.GenerateRandomString(20);
            var postResponse = await _fixture.Client.PostAsJsonAsync($"api/users/{user.Id}/new_workspace", name);

            try
            {
                postResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var workspace = await postResponse.Content.ReadFromJsonAsync<Workspace>();
            Assert.NotNull(workspace);
            Assert.Equal(name, workspace.Name);          
        }
    }

    [Fact]
    public async Task TestDeleteWorkspaceByEntity()
    {
        foreach (var entity in _fixture.DeleteEntitiesWithWorkspaces)
        {
            var getResponse = await _fixture.Client.GetAsync($"api/users/{entity}/initial");

            try
            {
                getResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var user = await getResponse.Content.ReadFromJsonAsync<User>();
            Assert.NotNull(user);

            foreach (var workspace in user.Workspaces)
            {
                var deleteResponse = await _fixture.Client.DeleteAsync($"api/users/{user.Id}/workspace/{workspace.Id}");

                try
                {
                    deleteResponse.EnsureSuccessStatusCode();
                }
                catch (HttpRequestException exception)
                {
                    Assert.Fail(exception.Message);
                }
            }
        }
    }

    [Fact]
    public async Task TestAddChatByEntity()
    {
        foreach (var entity in _fixture.EntitiesWithWorkspaces)
        {
            var getResponse = await _fixture.Client.GetAsync($"api/users/{entity}/initial");

            try
            {
                getResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var user = await getResponse.Content.ReadFromJsonAsync<User>();
            Assert.NotNull(user);

            foreach (var workspace in user.Workspaces)
            {
                var topic = RandomString.GenerateRandomString(20);
                var postResponse = await _fixture.Client.PostAsJsonAsync($"api/users/{user.Id}/workspace/{workspace.Id}/new_chat", topic);

                try
                {
                    postResponse.EnsureSuccessStatusCode();
                }
                catch (HttpRequestException exception)
                {
                    Assert.Fail(exception.Message);
                }

                var chat = await postResponse.Content.ReadFromJsonAsync<ChatRoom>();
                Assert.NotNull(chat);
                Assert.Equal(topic, chat.Topic);
            }
        }
    }

    [Fact]
    public async Task TestDeleteChatByEntity()
    {
        foreach (var entity in _fixture.DeleteEntitiesWithChatRooms)
        {
            var getResponse = await _fixture.Client.GetAsync($"api/users/{entity}/initial");

            try
            {
                getResponse.EnsureSuccessStatusCode();
            }
            catch (HttpRequestException exception)
            {
                Assert.Fail(exception.Message);
            }

            var user = await getResponse.Content.ReadFromJsonAsync<User>();
            Assert.NotNull(user);

            foreach (var workspace in user.Workspaces)
            {
                foreach (var chatroom in workspace.ChatRooms)
                {
                    var deleteResponse = await _fixture.Client.DeleteAsync($"api/users/{user.Id}/workspace/{workspace.Id}/chat/{chatroom.Id}");

                    try
                    {
                        deleteResponse.EnsureSuccessStatusCode();
                    }
                    catch (HttpRequestException exception)
                    {
                        Assert.Fail(exception.Message);
                    }
                }
            }
        }
    }
}

