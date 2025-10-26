using System.Net.Http.Headers;
using System.Net.Http.Json;
using UserManagementSystem.Models;

namespace MyWebApi.Tests;

[Collection("S3Collection")]
public class TestImageContext(ImageEndpointFixture fixture) : IClassFixture<ImageEndpointFixture>
{
    private readonly ImageEndpointFixture _fixture = fixture;

    [Fact]
    public async Task GetImage()
    {
        var getURLRequest = new HttpRequestMessage(HttpMethod.Get, $"api/image/{_fixture.GetRequestID}");
        var getURLResponse = await _fixture.Client.SendAsync(getURLRequest);
        Assert.NotNull(getURLResponse);

        try
        {
            getURLResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail($"{getURLRequest.RequestUri}{exception.Message}");
        }

        using var client = new HttpClient();

        var getCredentials = await getURLResponse.Content.ReadFromJsonAsync<GetImageModel>();
        Assert.NotNull(getCredentials);

        var getImageRequest = new HttpRequestMessage(HttpMethod.Get, getCredentials.URL);
        var getImageResponse = await client.SendAsync(getImageRequest);
        Assert.NotNull(getImageRequest);

        try
        {
            getImageResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail(exception.Message);
        }

        Assert.NotNull(getImageResponse.Content);
    }

    [Fact]
    public async Task GetImageFailed()
    {
        var getURLRequest = new HttpRequestMessage(HttpMethod.Get, $"api/image/{Guid.NewGuid()}.jpg");
        var getURLResponse = await _fixture.Client.SendAsync(getURLRequest);
        Assert.NotNull(getURLResponse);

        try
        {
            getURLResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail($"{getURLRequest.RequestUri}{exception.Message}");
        }

        var getCredentials = await getURLResponse.Content.ReadFromJsonAsync<GetImageModel>();
        Assert.NotNull(getCredentials);

        using var client = new HttpClient();

        var getImageRequest = new HttpRequestMessage(HttpMethod.Get, getCredentials.URL);
        var getImageResponse = await client.SendAsync(getImageRequest);
        Assert.NotNull(getImageRequest);

        Assert.Throws<HttpRequestException>(getImageResponse.EnsureSuccessStatusCode);
    }

    [Fact]
    public async Task UploadImage()
    {
        var postURLRequest = new HttpRequestMessage(HttpMethod.Post, "api/image/new")
        {
            Content = JsonContent.Create(new ImageUploadRequest { FileName = _fixture.FileName(), ContentType = "image/jpeg" })
        };
        postURLRequest.Content.Headers.ContentType = new MediaTypeHeaderValue("application/json");
        var postURLResponse = await _fixture.Client.SendAsync(postURLRequest);
        Assert.NotNull(postURLResponse);

        try
        {
            postURLResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail(exception.Message);
        }

        using var client = new HttpClient();
        using var stream = File.OpenRead(_fixture.FilePath());

        var postCredentials = await postURLResponse.Content.ReadFromJsonAsync<PostImageModel>();
        Assert.NotNull(postCredentials);
        
        var postImageRequest = new HttpRequestMessage(HttpMethod.Put, postCredentials.URL)
        {
            Content = new StreamContent(stream)
        };
        postImageRequest.Content.Headers.ContentType = new MediaTypeHeaderValue("image/jpeg");
        var postImageResponse = await client.SendAsync(postImageRequest);

        try
        {
            postImageResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail(exception.Message);
        }

        var getURLRequest = new HttpRequestMessage(HttpMethod.Get, $"api/image/{postCredentials.ID}");
        var getURLResponse = await _fixture.Client.SendAsync(getURLRequest);
        Assert.NotNull(getURLResponse);

        try
        {
            getURLResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail($"{getURLRequest.RequestUri}{exception.Message}");
        }

        var getCredentials = await getURLResponse.Content.ReadFromJsonAsync<GetImageModel>();
        Assert.NotNull(getCredentials);

        var getImageRequest = new HttpRequestMessage(HttpMethod.Get, getCredentials.URL);
        var getImageResponse = await client.SendAsync(getImageRequest);
        Assert.NotNull(getImageRequest);

        try
        {
            getImageResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail(exception.Message);
        }

        Assert.NotNull(getImageResponse.Content);
    }

    [Fact]
    public async Task DeleteImage()
    {
        var postURLRequest = new HttpRequestMessage(HttpMethod.Post, "api/image/new")
        {
            Content = JsonContent.Create(new ImageUploadRequest { FileName = _fixture.FileName(), ContentType = "image/jpeg" })
        };
        postURLRequest.Content.Headers.ContentType = new MediaTypeHeaderValue("application/json");
        var postURLResponse = await _fixture.Client.SendAsync(postURLRequest);
        Assert.NotNull(postURLResponse);

        try
        {
            postURLResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail(exception.Message);
        }

        using var client = new HttpClient();
        using var stream = File.OpenRead(_fixture.FilePath());

        var postCredentials = await postURLResponse.Content.ReadFromJsonAsync<PostImageModel>();
        Assert.NotNull(postCredentials);
        
        var postImageRequest = new HttpRequestMessage(HttpMethod.Put, postCredentials.URL)
        {
            Content = new StreamContent(stream)
        };
        postImageRequest.Content.Headers.ContentType = new MediaTypeHeaderValue("image/jpeg");
        var postImageResponse = await client.SendAsync(postImageRequest);

        try
        {
            postImageResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail(exception.Message);
        }

        var deleteRequest = new HttpRequestMessage(HttpMethod.Delete, $"api/image/{postCredentials.ID}/delete");
        var deleteResponse = await _fixture.Client.SendAsync(deleteRequest);
        Assert.NotNull(deleteResponse);

        try
        {
            deleteResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail($"{deleteRequest.RequestUri}{exception.Message}");
        }

        var getURLRequest = new HttpRequestMessage(HttpMethod.Get, $"api/image/{postCredentials.ID}");
        var getURLResponse = await _fixture.Client.SendAsync(getURLRequest);
        Assert.NotNull(getURLResponse);

        try
        {
            getURLResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail(exception.Message);
        }

        var getCredentials = await getURLResponse.Content.ReadFromJsonAsync<GetImageModel>();
        Assert.NotNull(getCredentials);

        var getImageRequest = new HttpRequestMessage(HttpMethod.Get, getCredentials.URL);
        var getImageResponse = await client.SendAsync(getImageRequest);
        Assert.NotNull(getImageRequest);

        Assert.Throws<HttpRequestException>(getImageResponse.EnsureSuccessStatusCode);
    }
}