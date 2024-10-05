using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using Microsoft.IdentityModel.Tokens;
using UserManager.PublicKey;
using UserManager.Redis.Data;

namespace UserManager.Middlewares;

public class AuthMiddleware(KeyManager key, IServiceScopeFactory scopeFactory) : IMiddleware
{
    private readonly RsaSecurityKey _key = key.Get;
    private readonly IServiceScopeFactory _scopeFactory = scopeFactory;

    public async Task InvokeAsync(HttpContext context, RequestDelegate next)
    {
        var bearer = context.Request.Headers.Authorization.ToString();
        if (bearer != null)
        {
            var claims = ValidateToken(bearer);
            
            if (claims != null)
            {
                context.Items["entity"] = VerifyUser(claims);
            }
        }

        await next(context);   
    } 

    private ClaimsPrincipal? ValidateToken(string header)
    {
        var token = header["Bearer ".Length..].Trim();
        
        var tokenHandler = new JwtSecurityTokenHandler();
        try {
            var validationParams = new TokenValidationParameters
            {
                ValidateIssuerSigningKey = true,
                IssuerSigningKey = _key,
                ValidateIssuer = false,
                ValidateAudience = false,
                ValidateLifetime = false
            };

            return tokenHandler.ValidateToken(token, validationParams, out _);
        }
        catch {
            return null;
        }
    }

    private async Task<string?> VerifyUser(ClaimsPrincipal claims)
    {
        using var scope = _scopeFactory.CreateAsyncScope();
        var context = scope.ServiceProvider.GetRequiredService<ISessionRepo>();

        foreach (var claim in claims.Claims)
        {
            Console.WriteLine($"Type: {claim.Type}\nValue: {claim.Value}\n");
        }

        var accountEntity = claims.Claims.FirstOrDefault(claim => claim.Type == "account")?.Value;
        if (accountEntity == null)
        {
            return null;
        }

        var session = await context.GetSession(accountEntity);
        if (session == null)
        {
            return null;
        }

        var aboType = claims.Claims.FirstOrDefault(claim => claim.Type == "abo")?.Value;
        if (aboType == null)
        {
            return null;
        }
        if (!IsValidAboType(aboType))
        {
            return null;
        }

        var name = claims.Claims.FirstOrDefault(claim => claim.Type == "name")?.Value;
        if (name == null)
        {
            return null;
        }
        else if(!name.Equals(session.Name))
        {
            return null;
        }

        var sessionID = claims.Claims.FirstOrDefault(claim => claim.Type == "session")?.Value;
        if (sessionID == null)
        {
            return null;
        }
        else if (!sessionID.Equals(session.Id))
        {
            return null;
        }

        return accountEntity;
    }

    private static bool IsValidAboType(string aboType)
    {
        return aboType switch
        {
            "COMMUNITY_USER" => true,
            "PRO_USER" => true,
            "NONE" => false,
            _ => false,
        };
    }
}

