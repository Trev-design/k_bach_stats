namespace UserManagementSystem.Models;

public class ImageModel
{
    public string ID { get; init; } = string.Empty;
    public string URL { get; init; } = string.Empty;
} 

public class GetImageModel : ImageModel { }
public class PostImageModel : ImageModel { }
public class DeleteImageModel : ImageModel { }