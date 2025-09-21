namespace MyWebApi.Tests;

[Collection("Database collection")]
public class MigrationTestProfileController(EndpointsFixture fixture)
{
    private readonly EndpointsFixture _fixture = fixture;

    [Fact]
    public async Task TestGet()
    {

    }

    [Fact]
    public async Task TestGetFailed()
    {

    }

    [Fact]
    public async Task TestChangeImage()
    {

    }

    [Fact]
    public async Task TestChangeImageFailed()
    {

    }

    [Fact]
    public async Task TestChangeDescription()
    {

    }

    [Fact]
    public async Task TestChangeDescriptionFailed()
    {
        
    }
}