from sqlalchemy import create_engine, Column, Integer, String, DateTime, ForeignKey
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from datetime import datetime

# Database setup
engine = create_engine("sqlite:///bot.db")  # SQLite for local development
Session = sessionmaker(bind=engine)
session = Session()
Base = declarative_base()

# Models
class User(Base):
    __tablename__ = "users"
    id = Column(Integer, primary_key=True)  # Discord user ID
    song_count = Column(Integer, default=0)

class Song(Base):
    __tablename__ = "songs"
    id = Column(Integer, primary_key=True)
    song_name = Column(String, nullable=False)  # Song name (string)
    requested_by = Column(Integer, ForeignKey("users.id"))

# Create tables
Base.metadata.create_all(engine)

