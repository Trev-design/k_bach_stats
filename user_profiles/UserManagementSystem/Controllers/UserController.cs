using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using UserManagementSystem.Models;
using UserManagementSystem.Services.Database;

namespace UserManagementSystem.Controllers;

[Controller]
public class UserController(AppDBContext dbContext) : Controller
{
    private readonly AppDBContext _dbContext = dbContext;

    [HttpGet("api/users/{entity}/initial")]
    public async Task<ActionResult<User>> GetInitial(string entity)
    {
        var user = await UserDBImpl.GetWholeUser(_dbContext, entity);
        if (user == null) return NotFound();
        return Ok(user);
    }

    [HttpGet("api/users/{id}")]
    public async Task<ActionResult<User>> Get(Guid id)
    {
        var user = await UserDBImpl.GetUserById(_dbContext, id);
        if (user == null) return NotFound();
        return Ok(user);
    }

    [HttpPost("api/users/{id}/new_workspace")]
    public async Task<ActionResult<Workspace>> NewWorkspace(Guid id, [FromBody] string workspaceId)
    {
        var workspace = await UserDBImpl.AddNewWorkspace(_dbContext, id, workspaceId);
        if (workspace == null) return NotFound();
        return Ok(workspace);
    }

    [HttpDelete("api/users/{id}/workspace/{workspaceId}")]
    public async Task<ActionResult> DeleteWorkspace(Guid id, Guid workspaceId)
    {
        try
        {
            await UserDBImpl.DeleteWorkspace(_dbContext, id, workspaceId);
        }
        catch
        {
            return NotFound();
        }
        return NoContent();
    }

    [HttpPost("api/users/{id}/workspace/{workspaceId}/new_chat")]
    public async Task<ActionResult<ChatRoom>> NewChat(Guid id, Guid workspaceId, [FromBody] string topic)
    {
        var chatRoom = await UserDBImpl.NewChatRoom(_dbContext, id, workspaceId, "dummyref", topic);
        if (chatRoom == null) return NotFound();
        return Ok(chatRoom); 
    }

    [HttpDelete("api/users/{id}/workspace/{workspaceId}/chat/{chatId}")]
    public async Task<ActionResult> DeleteChat(Guid id, Guid workspaceId, Guid chatId)
    {
        try
        {
            await UserDBImpl.DeleteChat(_dbContext, id, workspaceId, chatId);
        }
        catch
        {
            return NotFound();
        }
        return NoContent();
    }

    [HttpDelete("api/users/{id}")]
    public async Task<ActionResult> Delete(Guid id)
    {
        try
        {
            await UserDBImpl.DeleteUser(_dbContext, id);
        }
        catch
        {
            return NotFound();
        }
        return NoContent();
    }
}

