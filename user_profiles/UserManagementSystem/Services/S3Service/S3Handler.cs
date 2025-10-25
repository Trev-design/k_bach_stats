using Amazon.S3;
using Amazon.S3.Model;
using Minio;
using Minio.DataModel.Args;
using Org.BouncyCastle.Asn1.Cms;
using UserManagementSystem.Models;

namespace UserManagementSystem.Services.S3Service;

public class S3Handler(IMinioClient client, S3Settings settings)
{
    private readonly IMinioClient _client = client;
    private readonly S3Settings _settings = settings;

    public async Task<GetImageModel?> GetImageCredentials(string id)
    {
        try
        {
            var url = await _client.PresignedGetObjectAsync(new PresignedGetObjectArgs()
                .WithBucket(_settings.BucketName)
                .WithObject(id)
                .WithExpiry(60 * 10));

            return new GetImageModel { ID = id, URL = url };
        }
        catch
        {
            return null;
        }
    }

    public async Task<PostImageModel?> PostImageCredentials(string fileName, string contentType)
    {
        try
        {
            string id = $"images/{Guid.NewGuid()}{Path.GetExtension(fileName)}";

            var url = await _client.PresignedPutObjectAsync(new PresignedPutObjectArgs()
            .WithBucket(_settings.BucketName)
            .WithObject(id)
            .WithExpiry(10 * 60)
            .WithHeaders(new Dictionary<string, string>
            {
                { "x-amz-meta-original-filename", fileName }
            }));

            return new PostImageModel { ID = id, URL = url };
        } catch
        {
            return null;
        }
    }

    public async Task DeleteImageRequest(string id)
    {
        await _client.RemoveObjectAsync(new RemoveObjectArgs()
        .WithBucket(_settings.BucketName)
        .WithObject(id));
    }
}
