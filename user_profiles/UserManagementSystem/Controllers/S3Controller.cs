using Microsoft.AspNetCore.Mvc;
using UserManagementSystem.Models;
using UserManagementSystem.Services.S3Service;

namespace UserManagementSystem.Controllers;

[Controller]
[Route("api/[controller]")]
public class ImageController(S3Handler handler) : Controller
{
    private readonly S3Handler _handler = handler; 

    [HttpGet("{id}")]
    public async Task<ActionResult<GetImageModel>> Get(string id)
    {
        var result = await _handler.GetImageCredentials(id);
        if (result == null) return BadRequest("something went wrong");
        return Ok(result);
    }

    [HttpPost("new")]
    public async Task<ActionResult<PostImageModel>> Post([FromBody] ImageUploadRequest request)
    {
        var result = await _handler.PostImageCredentials(request.FileName, request.ContentType);
        if (result == null) return BadRequest("something went wrong");
        return Ok(result);
    }

    [HttpDelete("{id}/delete")]
    public async Task<ActionResult<string>> Delete(string id)
    {
        await _handler.DeleteImageRequest(id);
        return Ok("image deleted");
    }
}