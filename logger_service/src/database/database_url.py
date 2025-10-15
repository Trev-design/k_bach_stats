import os


def get_url():
  user = os.environ.get("DATABASE_USER", "test_user")
  password = os.environ.get("DATABASE_PASSWORD", "test_password")
  host = os.environ.get("DATABASE_HOST", "localhost")
  port = os.environ.get("DATABASE_PORT", "5432")
  database = os.environ.get("DATABASE_NAME", "test_db")

  return f"postgresql+asyncpg://{user}:{password}@{host}:{port}/{database}"
