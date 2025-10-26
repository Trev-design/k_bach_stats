
using System.Net.Http.Headers;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.TestHost;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using UserManagementSystem.Controllers;
using UserManagementSystem.Services.S3Service;

namespace MyWebApi.Tests;

public class ImageEndpointFixture(S3Fixture s3Fixture) : IAsyncLifetime
{
    private readonly S3Fixture _s3Fixture = s3Fixture;
    private IHost _host = null!;
    public HttpClient Client { get; private set; } = null!;
    public string GetRequestID { get; private set; } = null!;

    public string FilePath() => _s3Fixture.FilePath;
    public string FileName() => _s3Fixture.FileName;

    public async Task DisposeAsync()
    {
        Client.Dispose();
        await _host.StopAsync();
        _host.Dispose();
    }

    public async Task InitializeAsync()
    {
        var builder = WebApplication.CreateBuilder();

        builder.Services.AddSingleton(_s3Fixture.Settings);
        builder.Services.AddSingleton(_s3Fixture.Client);
        builder.Services.AddSingleton<S3Handler>();

        builder.Services.AddControllers().AddApplicationPart(typeof(ImageController).Assembly);
        builder.WebHost.UseTestServer();

        var app = builder.Build();

        app.MapControllers();

        await app.StartAsync();

        Client = app.GetTestClient();

        await UploadGetImagae();

        _host = app;
    }

    private async Task UploadGetImagae()
    {
        var postCreds = await _s3Fixture.Handler.PostImageCredentials(_s3Fixture.FileName, "image/jpeg");
        Assert.NotNull(postCreds);

        using var client = new HttpClient();
        using var stream = File.OpenRead(_s3Fixture.FilePath);

        var postRequest = new HttpRequestMessage(HttpMethod.Put, postCreds.URL)
        {
            Content = new StreamContent(stream)
        };

        postRequest.Content.Headers.ContentType = new MediaTypeHeaderValue("image/jpeg");
        //postRequest.Headers.Add("x-amz-meta-original-filename", "test_image.jpg");

        await client.SendAsync(postRequest);

        GetRequestID = postCreds.ID;
    }
}