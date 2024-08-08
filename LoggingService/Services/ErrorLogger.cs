using System.Diagnostics.Metrics;
using Grpc.Core;
using LoggingService.Data;
using LoggingService.Models;
using LogMessage;

namespace LoggingService.Services;

public class ErrorLogger : LogMessageService.LogMessageServiceBase
{
    private readonly AppDBContext appDBContext;

    public ErrorLogger(AppDBContext dbContext)
    {
        appDBContext = dbContext;
    }

    public override async Task<LogMessageResponse> SendLogMessage(IAsyncStreamReader<LogMessageRequest> requestStream, ServerCallContext context)
    {
        int count = 0;

        await foreach(var message in requestStream.ReadAllAsync())
        {
            if (message.Header == string.Empty || message.Error == string.Empty || message.Log == string.Empty)
            {
                count++;
            }
            else
            {
                var logMessage = new Log {
                    Id = Guid.NewGuid(),
                    DescriptionHeader = message.Header,
                    LogMessage = message.Log,
                    Error = message.Error
                };

                await appDBContext.AddAsync(logMessage);
                await appDBContext.SaveChangesAsync();
            }
        }

        return await Task.FromResult(new LogMessageResponse {
            Result = count > 0 ? string.Format("There were {0} invalid credentials", count) : "OK"
        });
    }
}