using DotNetEnv;

namespace UserManagementSystem.Utils;

public class DBConnstring
{
    public static string GetLocalEnvConnectionString()
    {
        Env.Load();

        string host = Env.GetString("MY_SQL_HOST");
        string user = Env.GetString("MY_SQL_USER");
        string pass = Env.GetString("MY_SQL_PASS");
        string name = Env.GetString("MY_SQL_NAME");

        return string.Format("Server={0};Database={1};User ID={2};Password={3};", host, name, user, pass);
    }

    public static string GetDockerEnvConnectionString()
    {
        string host = Environment.GetEnvironmentVariable("MY_SQL_HOST") ?? "localhost";
        string user = Environment.GetEnvironmentVariable("MY_SQL_USER") ?? throw new Exception("Must have a username");
        string pass = Environment.GetEnvironmentVariable("MY_SQL_PASS") ?? throw new Exception("Must have a password");
        string name = Environment.GetEnvironmentVariable("MY_SQL_NAME") ?? "mysql";

        return string.Format("Server={0};Database={1};User ID={2};Password={3};", host, name, user, pass);
    }
}