using Microsoft.AspNetCore.Mvc;
using UserManagementSystem.Models;
using UserManagementSystem.Services.Database;

namespace UserManagementSystem.Controllers;

[Controller]
[Route("api/[controller]")]
public class ProfileController(AppDBContext dbContext) : Controller
{
    private readonly AppDBContext _dbContext = dbContext;

    [HttpGet("{id}")]
    public async Task<ActionResult<Profile>> Get(Guid id)
    {
        var profile = await ProfileDBImpl.GetProfile(_dbContext, id);
        if (profile == null) return NotFound();
        return Ok(profile);
    }

    [HttpPut("{id}/new_image")]
    public async Task<ActionResult> ChangeImage(Guid id, [FromBody] string imagePath)
    {
        try
        {
            await ProfileDBImpl.ChangeImage(_dbContext, id, imagePath);
        }
        catch
        {
            return NotFound();
        }
        return NoContent();
    }

    [HttpPut("{id}/new_description")]
    public async Task<ActionResult> ChangeDescription(Guid id, [FromBody] string description)
    {
        try
        {
            await ProfileDBImpl.ChangeDescription(_dbContext, id, description);
        }
        catch
        {
            return NotFound();
        }
        return NoContent();
    }
}