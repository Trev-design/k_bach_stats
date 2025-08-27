using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using UserManagementSystem.Models;
using UserManagementSystem.Services;

namespace UserManagementSystem.Controllers;

[Controller]
[Route("api/[controller]")]
public class UserController(AppDBContext dbContext) : Controller
{
    private readonly AppDBContext _dbContext = dbContext;

    [HttpGet("{id:guid}")]
    public async Task<ActionResult<User>> Get(Guid id)
    {
        var user = await UserDBImpl.GetUserById(_dbContext, id);
        if (user == null) return NotFound();
        return user;
    }

    [HttpPost("{id:guid}/new_workspace")]
    public async Task<ActionResult> NewWorkspace(Guid id, [FromBody] string workspaceId)
    {
        var workspace = await UserDBImpl.AddNewWorkspace(_dbContext, id, workspaceId);
        if (workspace == null) return NotFound();
        return Ok(workspace);
    }

    [HttpDelete("{id:guid}/workspace/{workspaceId:guid}")]
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

    [HttpPost("{id:guid}/workspace/{workspaceId:guid}/new_chat")]
    public async Task<ActionResult> NewChat(Guid id, Guid workspaceId, [FromBody] string topic)
    {
        var chatRoom = await UserDBImpl.NewChatRoom(_dbContext, id, workspaceId, "dummyref", topic);
        if (chatRoom == null) return NoContent();
        return Ok(chatRoom); 
    }

    [HttpDelete("{id:guid}/workspace/{workspaceId:guid}/chat/{chatId:guid}")]
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

    [HttpDelete("{id:guid}")]
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