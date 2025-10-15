from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.sql import func
from sqlalchemy import DateTime
from datetime import datetime
import uuid


class Base(DeclarativeBase):
  pass


class Error(Base):
  __tablename__ = "errors"
  id: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
  message: Mapped[str] = mapped_column(nullable=False)
  service: Mapped[str] = mapped_column(nullable=False)
  file_name: Mapped[str] = mapped_column(nullable=False)
  line: Mapped[int] = mapped_column(nullable=False)
  happened_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now())


class Event(Base):
  __tablename__ = "events"
  id: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
  action: Mapped[str] = mapped_column(nullable=False)
  workspace_name: Mapped[str] = mapped_column(nullable=False)
  entity: Mapped[str] = mapped_column(nullable=False)
  happened_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now())


class WaitingLogs(Base): 
  __tablename__ = "waiting_logs"
  id: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
  payload: Mapped[str] = mapped_column(nullable=False)


