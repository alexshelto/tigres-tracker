import discord

from discord.ext import commands
from sqlalchemy import create_engine, func
from sqlalchemy.orm import sessionmaker
from config import config
from db.models import Song, User
import re
from dataclasses import dataclass

# DB Setup
engine = create_engine("sqlite:///bot.db")
Session = sessionmaker(bind=engine)
session = Session()


intents = discord.Intents.default()
intents.message_content = True #Required to read message content
intents.members = True

client = commands.Bot(command_prefix="t!", intents=intents)

@client.event
async def on_ready():
    print(f'bot is now online and ready')


@client.event
async def on_message(message):
    if message.author == client.user:
        return


    # Check for embedded message
    if message.embeds:
        for embed in message.embeds:
            embed_data = embed.to_dict()
            song_info = process_embed_data_for_now_playing(embed_data)
            if song_info is not None:
                print(song_info)
                song = Song(
                        song_name=song_info.name, 
                        requested_by=song_info.requested_by
                    )
                session.add(song)

                # increment_song_count(song_info.requested_by)
                session.commit()
                print(f'Saved song: {song.song_name} by {song.requested_by}')

@dataclass
class ParsedSongInfo:
    name: str
    requested_by: int


def process_embed_data_for_now_playing(data: dict[str,str]) -> ParsedSongInfo | None: 
    song_str = None
    requested_by = None

    title = data.get("title", "")
    description = data.get("description", "")
    
    if title.lower().strip() != "now playing": 
        return None

    description_lines = description.splitlines()
    song_str = description_lines[0]
    requested_by = description_lines[-1]
    requested_by_id = extract_user_id(requested_by)

    print(song_str)
    print(requested_by)

    if song_str and requested_by_id is not None and requested_by_id.isdigit():
        return ParsedSongInfo(song_str, int(requested_by_id))


def extract_user_id(request_string):
    match = re.search(r"<@(\d+)>", request_string)
    if match:
        return match.group(1)  # Return the user ID as a string
    else:
        return None  # If no match is found


def increment_song_count(requested_by_id: int):
    # Query the database to find the user
    user = session.query(User).filter_by(id=requested_by_id).first()
    
    if user:
        user.update({User.song_count: User.song_count + 1})
    else:
        # If the user is not found, create a new user record
        new_user = User(id=requested_by_id, song_count=1)
        session.add(new_user)

    try: 
        session.commit()
    except Exception as e:
        session.rollback()
        print(f'Error updating song count for user: {requested_by_id} {e}')





if config.BOT_TOKEN is None:
    print("Bot token is not set in config. Exiting")
    exit(1)

client.run(config.BOT_TOKEN)
