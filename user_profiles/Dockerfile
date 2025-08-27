# --- Build Stage ---
FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src

# Projektdateien kopieren und Restore
COPY /UserManagementSystem/*.csproj .
RUN dotnet restore

# Restlichen Code kopieren und Build
COPY . .
RUN dotnet publish -c Release -o /app/publish --no-restore

# --- Runtime Stage ---
FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS runtime
WORKDIR /app

COPY --from=build /app/publish .

EXPOSE 5132 7133
ENTRYPOINT ["dotnet", "UserManagementSystem.dll"]