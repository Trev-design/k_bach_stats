using System.Net.Http.Headers;
using System.Net.Http.Json;
using UserManagementSystem.Models;

namespace MyWebApi.Tests; 

[Collection("S3Collection")]
public class TestS3(S3Fixture fixture)
{
    private readonly S3Fixture _fixture = fixture;

    [Fact]
    public async Task TestUpload()
    {
        var postCreds = await _fixture.Handler.PostImageCredentials(_fixture.FileName, "image/jpeg");
        Assert.NotNull(postCreds);

        using var client = new HttpClient();
        using var stream = File.OpenRead(_fixture.FilePath);

        var postRequest = new HttpRequestMessage(HttpMethod.Put, postCreds.URL)
        {
            Content = new StreamContent(stream)
        };

        postRequest.Content.Headers.ContentType = new MediaTypeHeaderValue("image/jpeg");
        //postRequest.Headers.Add("x-amz-meta-original-filename", "test_image.jpg");

        var postResponse = await client.SendAsync(postRequest);
        Assert.NotNull(postResponse);

        try
        {
            postResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail($"{postCreds.URL} ERROR:{exception}");
        }

        var getCreds = await _fixture.Handler.GetImageCredentials(postCreds.ID);
        Assert.NotNull(getCreds);

        var getRequest = new HttpRequestMessage(HttpMethod.Get, getCreds.URL);

        var getResponse = await client.SendAsync(getRequest);

        try
        {
            getResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail(exception.Message);
        }

        Assert.NotNull(getResponse.Content);
    }

    [Fact]
    public async Task TestGetFailed()
    {
        var creds = await _fixture.Handler.GetImageCredentials($"{Guid.NewGuid()}.txt");
        Assert.NotNull(creds);

        using var client = new HttpClient();

        var request = new HttpRequestMessage(HttpMethod.Get, creds.URL);
        var response = await client.SendAsync(request);

        Assert.Throws<HttpRequestException>(() => { response.EnsureSuccessStatusCode(); });
    }

    [Fact]
    public async Task TestDelete()
    {
        var postCreds = await _fixture.Handler.PostImageCredentials(_fixture.FileName, "image/jpeg");
        Assert.NotNull(postCreds);

        using var client = new HttpClient();
        using var stream = File.OpenRead(_fixture.FilePath);

        var postRequest = new HttpRequestMessage(HttpMethod.Put, postCreds.URL)
        {
            Content = new StreamContent(stream)
        };

        postRequest.Content.Headers.ContentType = new MediaTypeHeaderValue("image/jpeg");

        var postResponse = await client.SendAsync(postRequest);
        Assert.NotNull(postResponse);

        try
        {
            postResponse.EnsureSuccessStatusCode();
        }
        catch (HttpRequestException exception)
        {
            Assert.Fail($"{postCreds.URL} ERROR:{exception}");
        }

        await _fixture.Handler.DeleteImageRequest(postCreds.ID);

        var getCreds = await _fixture.Handler.GetImageCredentials(postCreds.ID);
        Assert.NotNull(getCreds);

        var request = new HttpRequestMessage(HttpMethod.Get, getCreds.URL);
        var response = await client.SendAsync(request);

        Assert.Throws<HttpRequestException>(() => { response.EnsureSuccessStatusCode(); });
    }
} 