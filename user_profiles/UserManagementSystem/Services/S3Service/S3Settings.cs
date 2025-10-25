using DotNetEnv;

namespace UserManagementSystem.Services.S3Service;

public sealed class S3Settings
{
    public string RegionName { get; private set; } = string.Empty;
    public string BucketName { get; private set; } = string.Empty;
    public string URL { get; private set; }

    public S3Settings()
    {
        RegionName = Env.GetString("S3_REGION", "eu-central-1");
        BucketName = Env.GetString("S3_BUCKET", "userprofiles");

        URL = Env.GetString("S3_URL", "http://localhost:9000");
    }

    public string GetUser()
    {
        return Env.GetString("S3_USER", "tes_tuser");
    }

    public string GetSekretKey()
    {
        // condionaly use hashicorp client in production
        return Env.GetString("S3_PASSWORD", "test_password");
    }
}