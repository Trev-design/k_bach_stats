using System.Security.Cryptography;
using Microsoft.IdentityModel.Tokens;

namespace UserManager.PublicKey;

public class KeyManager 
{
    private readonly RsaSecurityKey _key;
    public KeyManager()
    {
        var pemKey = File.ReadAllText("public.pem") ?? throw new ArgumentException("could not convert pem to key");
        using var rsa = RSA.Create();
        rsa.ImportFromPem(pemKey.ToCharArray());

        _key = new RsaSecurityKey(rsa);
    }

    public RsaSecurityKey Get => _key;
}