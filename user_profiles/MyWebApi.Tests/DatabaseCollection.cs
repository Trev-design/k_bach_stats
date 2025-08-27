using Xunit;

namespace MyWebApi.Tests;

[CollectionDefinition("Database collection")]
public class DatabaseCollection : ICollectionFixture<DatabaseFixture>
{
    // Leere Klasse – dient nur der Gruppierung
}