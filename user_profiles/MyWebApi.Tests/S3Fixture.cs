using Minio;
using Minio.DataModel.Args;
using Testcontainers.Minio;
using UserManagementSystem.Services.S3Service;

namespace MyWebApi.Tests;

public class S3Fixture : IAsyncLifetime
{
    private MinioContainer _container = null!;
    private IMinioClient _client = null!;
    private S3Settings _settings = null!;
    public S3Handler Handler { get; private set; } = null!;
    public string FileName { get; } = "test_image.jpg";
    public string FilePath { get; private set; } = null!;

    public async Task DisposeAsync()
    {
        _client.Dispose();
        await _container.StopAsync();
        await _container.DisposeAsync();
    }

    public async Task InitializeAsync()
    {
        _settings = new S3Settings();

        _container = new MinioBuilder()
            .WithImage("minio/minio:RELEASE.2025-09-07T16-13-09Z-cpuv1")
            .WithEnvironment("MINIO_ROOT_USER", _settings.GetUser())
            .WithEnvironment("MINIO_ROOT_PASSWORD", _settings.GetSekretKey())
            .WithPortBinding(9000, true)
            .Build();

        await _container.StartAsync();

        var port = _container.GetMappedPublicPort();

        _client = new MinioClient()
            .WithEndpoint($"localhost:{port}")
            .WithCredentials(_settings.GetUser(), _settings.GetSekretKey())
            .WithSSL(false)
            .Build();

        bool found = await _client.BucketExistsAsync(new BucketExistsArgs().WithBucket(_settings.BucketName));
        if (!found) await _client.MakeBucketAsync(new MakeBucketArgs().WithBucket(_settings.BucketName));

        Handler = new S3Handler(_client, _settings);

        var tempFile = Path.Combine(Path.GetTempPath(), FileName);
        await File.WriteAllBytesAsync(tempFile, [0xFF, 0xD8, 0xFF, 0xD9]);
        FilePath = tempFile;
    }
}