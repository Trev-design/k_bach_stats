using Minio;
using Minio.DataModel.Args;
using Testcontainers.Minio;
using UserManagementSystem.Services.S3Service;

namespace MyWebApi.Tests;

public class S3Fixture : IAsyncLifetime
{
    private MinioContainer _container = null!;
    public IMinioClient Client { get; private set; } = null!;
    public S3Settings Settings { get; private set; } = null!;
    public S3Handler Handler { get; private set; } = null!;
    public string FileName { get; } = "test_image.jpg";
    public string FilePath { get; private set; } = null!;

    public async Task DisposeAsync()
    {
        Client.Dispose();
        await _container.StopAsync();
        await _container.DisposeAsync();
    }

    public async Task InitializeAsync()
    {
        Settings = new S3Settings();

        _container = new MinioBuilder()
            .WithImage("minio/minio:RELEASE.2025-09-07T16-13-09Z-cpuv1")
            .WithEnvironment("MINIO_ROOT_USER", Settings.GetUser())
            .WithEnvironment("MINIO_ROOT_PASSWORD", Settings.GetSekretKey())
            .WithPortBinding(9000, true)
            .Build();

        await _container.StartAsync();

        var port = _container.GetMappedPublicPort();

        Client = new MinioClient()
            .WithEndpoint($"localhost:{port}")
            .WithCredentials(Settings.GetUser(), Settings.GetSekretKey())
            .WithSSL(false)
            .Build();

        bool found = await Client.BucketExistsAsync(new BucketExistsArgs().WithBucket(Settings.BucketName));
        if (!found) await Client.MakeBucketAsync(new MakeBucketArgs().WithBucket(Settings.BucketName));

        Handler = new S3Handler(Client, Settings);

        var tempFile = Path.Combine(Path.GetTempPath(), FileName);
        await File.WriteAllBytesAsync(tempFile, [0xFF, 0xD8, 0xFF, 0xD9]);
        FilePath = tempFile;
    }
}