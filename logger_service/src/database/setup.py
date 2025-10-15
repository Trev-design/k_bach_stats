from .models import Base
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker
from .database_url import get_url

engine = create_async_engine(get_url(), connect_args={"check_same_thread": False})
session = async_sessionmaker(engine)


async def init_db():
  async with engine.begin() as connection:
    await connection.run_sync(Base.metadata.create_all)


async def get_db():    
  db = session()
  try:
    yield db
  finally:
    await db.close()