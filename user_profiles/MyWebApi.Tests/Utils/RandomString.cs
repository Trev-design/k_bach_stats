using System.Security.Cryptography;

namespace MyWebApi.Tests.Utils;

public class RandomString
{
    public static string GenerateRandomString(int maxLength)
    {
        var length = GenerateRandomStringLength(maxLength);
        var buffer = new byte[length];
        RandomNumberGenerator.Fill(buffer);
        var randomString = buffer.ToString() ?? throw new Exception();
        return randomString;
    }

    public static string GenerateRandomEmail(int maxLength)
    {
        var first = GenerateRandomString(maxLength);
        var second = GenerateRandomString(maxLength);
        return string.Format("{0}@{1}.com", first, second);
    }

    private static int GenerateRandomStringLength(int maxLength)
    {
        return RandomNumberGenerator.GetInt32(maxLength);
    }
}