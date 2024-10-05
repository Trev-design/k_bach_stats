namespace UserManager.Middlewares;

public class CorsMiddleware : IMiddleware
{
    public async Task InvokeAsync(HttpContext context, RequestDelegate next)
    {
        var origin = context.Request.Headers.Origin.ToString();

        Console.WriteLine($"the content of the origin: {origin}");

        if (!string.IsNullOrEmpty(origin))
        {
            Console.WriteLine($"the cool content of the origin: {origin}");
            var uri = new Uri(origin);

            if (uri.Host == "localhost" && uri.Port == 5173)
            {
                Console.WriteLine($"the very cool content of the origin {origin}");

                context.Response.Headers.Append("Access-Control-Allow-Origin", origin);
                context.Response.Headers.Append("Access-Control-Allow-Headers", "Authorization, Content-Type");
                context.Response.Headers.Append("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS");
                context.Response.Headers.Append("Access-Control-Allow-Credentials", "true");
            }
        }

        if (context.Request.Method == "OPTIONS")
        {
            context.Response.StatusCode = 204;
            return;
        }

        await next(context);
    }
}